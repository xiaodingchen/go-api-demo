package xtrace

import (
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
	"go.uber.org/zap"
	"io"
	"strconv"
	"test.local/pkg/utils"
	"test.local/pkg/xzap"
)

var TraceCloser io.Closer

func NewJaegerTracer(service, jaegerLogPath, jaegerLogSamplingRate, jaegerLogBufferSize string) (err error) {
	samplingRate, _ := strconv.ParseFloat(jaegerLogSamplingRate, 64)
	bufferSize, _ := strconv.Atoi(jaegerLogBufferSize)
	if jaegerLogPath != "" {
		sampler := jaeger.NewRateLimitingSampler(samplingRate)
		reporter, err := NewJaegerFileReport(jaegerLogPath, bufferSize)
		if err != nil {
			Error("set jaeger report err", err)
			return err
		}

		zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
		injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
		extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
		options := []jaeger.TracerOption{}
		options = append(options, injector, extractor)
		var jaegerTracer opentracing.Tracer
		jaegerTracer, TraceCloser = jaeger.NewTracer(service, sampler, reporter, options...)
		opentracing.SetGlobalTracer(jaegerTracer)
		return nil
	}

	return errors.New("jaeger tracer set err")
}

//func NewJaegerTracerWithConfig() {
//	cfg := jaegercfg.Configuration{
//		Sampler: &jaegercfg.SamplerConfig{
//			Type:  jaeger.SamplerTypeConst,
//			Param: 1,
//		},
//		Reporter: &jaegercfg.ReporterConfig{
//			LogSpans:           true,
//			LocalAgentHostPort: "{host}:6831", // 替换host
//		},
//	}
//
//	closer, err := cfg.InitGlobalTracer(
//		"serviceName",
//	)
//}

type jaegerLog struct {
	logger *zap.Logger
}

func newjaegerLog(jaegerLogPath string) *jaegerLog {
	cfg := &xzap.LogConfig{
		Level:       -1,
		Dev:         true,
		OutputPaths: []string{jaegerLogPath},
	}

	j := &jaegerLog{}

	logger, err := xzap.NewZap(cfg)
	if err != nil {
		utils.Log().Error("[jaeger] set logger err", zap.Error(err))
	}
	if logger != nil {
		j.logger = logger
	}

	return j
}

func (j *jaegerLog) Infof(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	if j.logger != nil {
		//defer j.logger.Sync()
		j.logger.Info(msg)
	}
}

func (j *jaegerLog) Error(msg string) {
	if j.logger != nil {
		//defer j.logger.Sync()
		j.logger.Error(msg)
	}
}

func jaegerLogger(jaegerLogPath string) jaeger.Logger {
	logger := newjaegerLog(jaegerLogPath)
	if logger.logger == nil {
		return jaeger.StdLogger
	}

	return logger
}

func CloseJaegerTracer() {
	if TraceCloser != nil {
		TraceCloser.Close()
	}
}
