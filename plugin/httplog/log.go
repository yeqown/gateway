// Package httplog to log http request
package httplog

import (
	"bytes"
	"net/http"
	"time"

	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/utils"
	log "github.com/yeqown/server-common/logger"
)

// New func: generate a new HTTPLogger
func New(logger *log.Logger) *HTTPLogger {
	return &HTTPLogger{
		logger: logger,
	}
}

// HTTPLogger ...
type HTTPLogger struct {
	logger *log.Logger
}

// Handle ...
func (h *HTTPLogger) Handle(ctx *plugin.Context) {
	rbw := &respBodyWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: ctx.ResponseWriter(),
	}
	ctx.SetResponseWriter(rbw)

	// timer
	start := time.Now()
	path := ctx.Path
	cpyReq := utils.CopyRequest(ctx.Request())
	fields := make(map[string]interface{})

	// continue process
	ctx.Next()

	end := time.Now()
	latency := end.Sub(start)
	clientIP := ctx.Request().RemoteAddr
	fields["requestForm"] = utils.ParseRequestForm(cpyReq) // set request
	fields["responseBody"] = rbw.body.String()             // set response

	// log
	h.logger.WithFields(fields).Infof("[Request] %v |%3d| %13v | %15s |%-7s %s",
		end.Format("2006/01/02 - 15:04:05"),
		rbw.status,
		latency,
		clientIP,
		ctx.Method,
		path,
	)
}

type respBodyWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w respBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w respBodyWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
