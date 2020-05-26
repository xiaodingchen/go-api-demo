package xtrace

import (
	"bufio"
	"encoding/json"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"test.local/pkg/utils"
	"time"
)

type JaegerFileReport struct {
	path       string
	handler    *os.File
	w          *bufio.Writer
	bufferSize int
}

func NewJaegerFileReport(reportFilePath string, size int) (jaeger.Reporter, error) {
	os.MkdirAll(filepath.Dir(reportFilePath), 0755)
	f, err := os.OpenFile(reportFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	r := &JaegerFileReport{
		handler:    f,
		path:       reportFilePath,
		bufferSize: size,
	}

	if r.bufferSize > 0 {
		r.w = bufio.NewWriterSize(r.handler, r.bufferSize)
	} else {
		r.w = bufio.NewWriter(r.handler)
	}

	go r.writeToFileWithBuffer()

	return r, nil
}

func (r *JaegerFileReport) Report(span *jaeger.Span) {
	data, err := json.Marshal(BuildSpan(span))
	if err != nil {
		Error("report err", err)
		return
	}

	data = append(data, '\n')
	_, err = r.w.Write(data)
	if err != nil {
		Error("report err", err)
	}
	return
	//Error("report err", r.w.WriteByte('\n'))
}

func (r *JaegerFileReport) Close() {
	r.sync()
	Error("Close err", r.handler.Close())
}

func (r *JaegerFileReport) writeToFileWithBuffer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.sync()
		}
	}
}

func (r *JaegerFileReport) sync() {
	Error("Flush err", r.w.Flush())
	Error("Sync err", r.handler.Sync())
}

func Error(msg string, err error) {
	if err == nil {
		return
	}
	utils.Log().Error("[jaeger] "+msg, zap.Error(err))
}
