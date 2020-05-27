package xtrace

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetSpanCtxFromCtx(ctx context.Context) context.Context{
	if ctx == nil{
		return context.Background()
	}

	if c, ok := ctx.(*gin.Context); ok{
		if c == nil || c.Request == nil{
			return context.Background()
		}

		return c.Request.Context()
	}

	return ctx
}

func OpStartSpan(name string, ctx context.Context, peerService string) (sp opentracing.Span){
	ctx = GetSpanCtxFromCtx(ctx)
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil{
		return
	}
	tracer := opentracing.GlobalTracer()
	sp = tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))

	if jaegerSpCtx, ok := sp.Context().(jaeger.SpanContext); ok{
		if !jaegerSpCtx.IsSampled(){
			sp.Finish()
			sp = nil
			return sp
		}
	}
	ext.SpanKindRPCClient.Set(sp)
	ext.PeerService.Set(sp, peerService)
	return
}

func redisStartSpan(name string, ctx context.Context, client *redis.Client)(sp opentracing.Span){
	ctx = GetSpanCtxFromCtx(ctx)
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil{
		return
	}
	tracer := opentracing.GlobalTracer()
	sp = tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
	if jaegerSpCtx, ok := sp.Context().(jaeger.SpanContext); ok{
		if !jaegerSpCtx.IsSampled(){
			sp.Finish()
			sp = nil
			return sp
		}
	}
	ext.SpanKindRPCClient.Set(sp)
	ext.PeerAddress.Set(sp, client.Options().Addr)
	ext.PeerService.Set(sp, "redis")
	return
}

// WrapperRedis 创建一个带有span的redis客户端
func WrapperRedis(ctx context.Context, client *redis.Client) *redis.Client{
	ctx = GetSpanCtxFromCtx(ctx)
	if ctx == nil{
		return client
	}
	clone := client.WithContext(ctx)
	clone.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			name := cmd.Name()
			sp := redisStartSpan("redis_" + name, ctx, client)
			if sp == nil{
				return oldProcess(cmd)
			}
			sp.SetTag("cmd", name)
			args := cmd.Args()
			if len(args) > 2{
				sp.SetTag("args", args[1].(string))
			}

			if len(args) >= 3 { // 多记录key(一般是第2个)参数后一个参数，最多记录32个字符，便于排查一些问题
				arg1, _ := args[2].(string)
				if len(arg1) > 32 {
					sp.SetTag("arg1", arg1[:32])
				} else {
					sp.SetTag("arg1", arg1)
				}
			}

			// 执行实际的redis命令
			err := oldProcess(cmd)
			if err != nil && err != redis.Nil{
				ext.Error.Set(sp, true)
				sp.LogKV("event", "error", "message", err.Error())
				sp.LogFields(tracelog.Error(err))
			}

			sp.Finish()
			return err
		}
	})
	
	clone.WrapProcessPipeline(func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error {
			return func(cmders []redis.Cmder) error {
				sp := redisStartSpan("redis_pipeline", ctx, client)
				if sp == nil{
					return oldProcess(cmders)
				}
				// 将命令和key拼接起来
				var cmdName, keys string
				for idx, cmd := range cmders {
					index := strconv.Itoa(idx)
					cmdName += index + ":" + cmd.Name() + " "
					args := cmd.Args()
					if len(args) >= 2 {
						keys += index + ":" + args[1].(string) + " "
					}
				}
				sp.SetTag("cmds", cmdName)
				sp.SetTag("args", keys)
				err := oldProcess(cmders) // 执行实际redis命令
				if err != nil {
					var errorBuf bytes.Buffer
					for idx, cmd := range cmders {
						if cmd.Err() != nil && cmd.Err().Error() != redis.Nil.Error() { // 出错了并且不是key不存在的错误
							errorBuf.Write([]byte(strconv.Itoa(idx) + ":" + cmd.Err().Error() + " "))
						}
					}
					errorStr := errorBuf.String()
					if errorStr != "" {
						ext.Error.Set(sp, true)
						sp.LogKV("event", "error", "message", errorStr)
					}
				}

				sp.Finish()
				return err
			}
	})

	return clone
}

// HttpRequestWithTrace curl加上trace
func HttpRequestWithTrace(ctx context.Context, request *http.Request, client *http.Client)(response *http.Response, err error){
	ctx = GetSpanCtxFromCtx(ctx)
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil{
		parentSpan = opentracing.StartSpan("http_client")
	}
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("http_client", opentracing.ChildOf(parentSpan.Context()))
	if jaegerSpanContext, ok := sp.Context().(jaeger.SpanContext); ok && !jaegerSpanContext.IsSampled() {
		sp.Finish() // return back to pool
		sp = nil
		goto SkipTrace
	}
	defer func() {
		if err != nil {
			ext.Error.Set(sp, true)
			sp.LogFields(tracelog.Error(err))
		}
		sp.Finish()
	}()
	ext.SpanKindRPCClient.Set(sp)
	ext.PeerAddress.Set(sp, fmt.Sprintf("%s%s", request.URL.Host, request.URL.Path))
	ext.PeerService.Set(sp, "http")
	sp.LogFields(tracelog.String("http.query", request.URL.RawQuery))
	if request.Body != nil {
		param, _ := ioutil.ReadAll(request.Body)
		sp.LogFields(tracelog.String("req_body", string(param)))
		request.Body = ioutil.NopCloser(bytes.NewReader(param))
	}
	_ = sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
SkipTrace:
	response, err = client.Do(request)
	if err != nil {
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	return
}