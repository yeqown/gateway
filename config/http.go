package config

import (
	"fmt"
	"log"
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
	h.GET("/plugins", api.PluginsGET)   //
	h.GET("/config", api.AllConfigsGET) // done

	// Proxy Path Rules
	h.GET("/plugin/proxy/config/pathrules", api.ProxyConfigPathsGET) // done
	h.GET("/plugin/proxy/config/pathrule/:id", api.ProxyConfigPathGET)
	h.POST("/plugin/proxy/config/pathrule", api.ProxyConfigPathPOST)
	h.PUT("/plugin/proxy/config/pathrule/:id", api.ProxyConfigPathPUT)
	h.DELETE("/plugin/proxy/config/pathrule/:id", api.ProxyConfigPathDELETE)

	// Proxy Server Rules
	h.GET("/plugin/proxy/config/srvrules", api.ProxyConfigSrvsGET)
	h.GET("/plugin/proxy/config/srvrule/:id", api.ProxyConfigSrvGET)
	h.POST("/plugin/proxy/config/srvrule", api.ProxyConfigSrvPOST)
	h.PUT("/plugin/proxy/config/srvrule/:id", api.ProxyConfigSrvPUT)
	h.DELETE("/plugin/proxy/config/srvrule/:id", api.ProxyConfigSrvDELETE)

	// Proxy ReverseServer
	//TODO: h.GET("/plugin/proxy/config/reversesrv", api.ProxyConfigReverseSrvGET)
	h.GET("/plugin/proxy/config/reversesrv/:group", api.ProxyConfigReverseSrvGroupGET)
	h.DELETE("/plugin/proxy/config/reversesrv/:group", api.ProxyConfigReverseSrvGroupDELETE)

	h.GET("/plugin/proxy/config/reversesrv/:group/:id", api.ProxyConfigReverseSrvGET)
	h.POST("/plugin/proxy/config/reversesrv/:group", api.ProxyConfigReverseSrvPOST)
	h.PUT("/plugin/proxy/config/reversesrv/:group/:id", api.ProxyConfigReverseSrvPUT)
	h.DELETE("/plugin/proxy/config/reversesrv/:group/:id", api.ProxyConfigReverseSrvDELETE)

	// Cache
	h.GET("/plugin/cache/configs", api.CacheConfigsGET)
	h.GET("/plugin/cache/config/:id", api.CacheConfigGET)
	h.POST("/plugin/cache/config", api.CacheConfigPOST)
	h.PUT("/plugin/cache/config/:id", api.CacheConfigPUT)
	h.DELETE("/plugin/cache/config/:id", api.CacheConfigDELETE)

	// Gate
	h.GET("/gate/config", api.GateConfigGET)
	h.PUT("/gate/config", api.GateConfigPUT)
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

	// request log
	form := utils.ParseRequestForm(utils.CopyRequest(req))
	logger.Logger.WithFields(map[string]interface{}{
		"form": form,
	}).Info("request with form")

	recoverHandle(handle)(w, req, params)
}

func recoverHandle(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		defer func() {
			if v := recover(); v != nil {
				errmsg := fmt.Sprintf("error: %v", v)
				logger.Logger.Error(errmsg)
				log.Printf("ConfigAPI.panic %s", debug.Stack())
				utils.ResponseJSON(w, code.NewCodeInfo(code.CodeSystemErr, errmsg))
			}
		}()

		h(w, req, params)
	}
}
