package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
    "github.com/uber/jaeger-client-go"
    "strconv"
)

func TraceLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        tracer := opentracing.GlobalTracer()
        operationName := c.Request.URL.Path
        var sp opentracing.Span
        // 从header中读取trace
        parentCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
        if parentCtx == nil {
            sp = tracer.StartSpan(operationName)
        } else {
            sp = tracer.StartSpan(operationName, opentracing.ChildOf(parentCtx))
        }

        ctx := opentracing.ContextWithSpan(c.Request.Context(), sp)
        c.Request = c.Request.WithContext(ctx)
        // 判断是否需要采样
        if jaegerSpanContext, ok := sp.Context().(jaeger.SpanContext); ok {
            if !jaegerSpanContext.IsSampled() {
                sp.Finish()
                c.Next()
                return
            }
        }

        tracer.Inject(sp.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
        ext.SpanKindRPCServer.Set(sp)
        ext.HTTPUrl.Set(sp, c.Request.RequestURI)
        ext.HTTPMethod.Set(sp, c.Request.Method)
        sp.SetTag("requestId", c.GetString(requestID))
        c.Next()
        status := c.Writer.Status()
        ext.HTTPStatusCode.Set(sp, uint16(status))
        if status > 400 {
            httperr := "httpCodeError:" + strconv.Itoa(status)
            ext.Error.Set(sp, true)
            sp.LogKV("event", "error", "message", httperr)
        }

        sp.Finish()
        return
    }
}