package config

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config/api"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

// New a ConfigAPI
func New() *HTTP {
	h := &HTTP{
		Prefix: "/gateapi",
		Router: httprouter.New(),
	}

	h.initRouter()

	return h
}

// HTTP to do some work with api
type HTTP struct {
	*httprouter.Router
	Prefix string
}

func (h *HTTP) initRouter() {
	h.GET("/plugins", api.PluginsGET)
	h.GET("/configs", api.AllConfigsGET)
}

type muxResponse struct {
	code.CodeInfo
}

// ServeHTTP serve request
func (h *HTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		resp = new(muxResponse)
	)

	req.URL.Path = strings.TrimPrefix(req.URL.Path, h.Prefix)
	handle, params, tsr := h.Lookup(req.Method, req.URL.Path)
	_, _ = params, tsr
	if handle == nil {
		logger.Logger.Infof("method: %s, path: %s", req.Method, req.URL.Path)
		code.FillCodeInfo(resp, code.NewCodeInfo(404, "ConfigAPI.Not Found"))
		utils.ResponseJSON(w, resp)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	recoverHandle(handle)(w, req, params)
}

func recoverHandle(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		defer func() {
			if v := recover(); v != nil {
				logger.Logger.Errorf("ConfigAPI.panic %s", debug.Stack())
			}
		}()

		h(w, req, params)
	}
}
