package utils

import (
	"go.uber.org/zap"
)

var logs map[string]*zap.Logger

const (
	DefaultLoggerName = "default"
	ApiLoggerName = "api"
	JaegerLoggerName = "jaeger"
)

func Log(name... string) *zap.Logger {
	if len(name) == 0{
		name = []string{DefaultLoggerName}
	}

	return logs[name[0]]
}

func SetLog(name string, logger *zap.Logger) {
	if logs == nil{
		logs = make(map[string]*zap.Logger)
	}

	logs[name] = logger
}

