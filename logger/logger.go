// Package logger define output to std or file
package logger

import (
	log "github.com/yeqown/server-common/logger"
)

var (
	// Logger is an internal logger for api-gateway to log
	Logger *log.Logger
)

// InitLogger call server-common to
func InitLogger(logPath string) (err error) {
	Logger, err = log.NewJSONLogger(logPath, "gateway.log", "debug")
	return err
}
