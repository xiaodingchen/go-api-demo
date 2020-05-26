package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"test.local/pkg/utils"
	"test.local/pkg/xzap"
	"time"
)

var LogFormatter = func(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	keys := param.Keys
	val, _ := json.Marshal(keys)
	jsonFormat := `{"level":"info","ts":"%s","msg":"[GIN] request info","requestId":"%v","status":%3d,"latency":"%v","client_ip":"%s","method":"%s","path":"%s", "uri":"%s", "error_msg":"%s","keys":%s,"size":%d}` + "\n"
	return fmt.Sprintf(jsonFormat,
		param.TimeStamp.Format("2006-01-02 15:04:05.000000"),
		keys[requestID],
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Request.URL.Path,
		param.Path,
		param.ErrorMessage,
		string(val),
		param.BodySize,
	)
}

var LoggerSkipPaths []string

func Logger(debug bool) gin.HandlerFunc {
	f := gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: LogFormatter,
		SkipPaths: LoggerSkipPaths,
	})

	if debug {
		return f
	}

	cfg := &xzap.LogConfig{
		Level:       -1,
		OutputPaths: viper.GetStringSlice("log.requestLog"),
		Dev:         true,
	}

	logger, err := xzap.NewZap(cfg)
	if err != nil {
		utils.Log().Error("[GIN] set logger err", zap.Error(err))
		return f
	}

	var skip map[string]struct{}
	if length := len(LoggerSkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range LoggerSkipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		keys := c.Keys
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		if _, ok := skip[path]; !ok {
			t := time.Now()
			if raw != "" {
				path = path + "?" + raw
			}
			//defer logger.Sync()
			logger.Info(
				"[GIN] request info",
				zap.Any("keys", keys),
				zap.String(requestID, c.GetString(requestID)),
				zap.Int("status", c.Writer.Status()),
				zap.Any("latency", t.Sub(start)),
				zap.String("client_ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("uri", path),
				zap.Int("size", c.Writer.Size()),
				zap.String("error_msg", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
		}
	}
}
