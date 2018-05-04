package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/server-common/code"
)

func response(w http.ResponseWriter, s string) {
	logger.Logger.Info(s)
	_, err := io.WriteString(w, s)
	if err != nil {
		logger.Logger.Errorf("response err: %s", err.Error())
	}
}

// ResponseString ...
func ResponseString(w http.ResponseWriter, s string) {
	response(w, s)
}

// ResponseJSON ...
func ResponseJSON(w http.ResponseWriter, i interface{}) {
	bs, err := json.Marshal(i)
	if err != nil {
		bs, _ = json.Marshal(code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		logger.Logger.Errorf("get an err: %s", err.Error())
	}

	// set header
	w.Header().Set("Content-Type", "application/json")

	response(w, string(bs))
}
