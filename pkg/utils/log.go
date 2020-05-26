package utils

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Log() *zap.Logger {
	return log
}

func SetLog(logger *zap.Logger) {
	log = logger
}

func RequestLogger() {

}
