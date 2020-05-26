package xzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"test.local/pkg/utils"
	"time"
)

type LogConfig struct {
	Level       int32
	OutputPaths []string
	Dev         bool

}

var loggers []*zap.Logger

func NewZap(name string, cfg *LogConfig) (log *zap.Logger, err error) {
	var (
		zcfg zap.Config
	)
	zcfg = zap.NewProductionConfig()
	zcfg.OutputPaths = cfg.OutputPaths
	zcfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zcfg.Development = cfg.Dev
	zcfg.EncoderConfig.EncodeTime = TimeEncoder
	log, err = zcfg.Build()
	if err != nil {
		return
	}
	newCore := zapcore.NewTee(
		log.Core(),
	)

	log = zap.New(newCore).WithOptions(zap.AddCaller())
	loggers = append(loggers, log)
	utils.SetLog(name, log)
	return
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(utils.DateFormatTimestamp))
}

func Sync() {
	sync()
}

func Async() {
	go func() {
		for {
			sync()
			time.Sleep(time.Second)
		}
	}()
}

func sync() {
	for _, logger := range loggers {
		if logger == nil {
			continue
		}
		logger.Sync()
	}
}
