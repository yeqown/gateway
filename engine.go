package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/yeqown/gateway/config"
	"github.com/yeqown/gateway/plugin"
	log "github.com/yeqown/server-common/logger"
)

// TIMEOUT string
const TIMEOUT = "timeout"

// Engine ...
type Engine struct {
	Plugins   []plugin.Plugin
	numPlugin int
	addr      string

	prefix string
	Logger *log.Logger // inner logger

	cfgAPI *config.HTTP // config api handler
}

func (e *Engine) init() {
	e.numPlugin = len(e.Plugins)
	e.prefix = "/gate"

	// init engine config api handler
	e.cfgAPI = config.New()
}

func (e *Engine) use(plgs ...plugin.Plugin) {
	e.Plugins = append(e.Plugins, plgs...)
	e.numPlugin += len(plgs)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path, e.prefix)
	// ctx := ctxPool.Get().(*plugin.Context)
	// defer ctxPool.Put(ctx)
	e.Logger.Info("new request recved")

	// generate a new context
	ctx := plugin.New(w, req, e.numPlugin, e.Plugins)

	// start call
	ctx.Next()

	// reset resource
	// ctx.Reset()
	return
}

// ListenAndServe ...
func (e *Engine) ListenAndServe(addr string) error {
	if addr != "" {
		e.addr = addr
	}
	e.init()

	mux := http.NewServeMux()
	mux.Handle(e.prefix+"/",
		http.TimeoutHandler(e, 5*time.Second, TIMEOUT))
	mux.Handle(e.cfgAPI.Prefix+"/",
		http.TimeoutHandler(e.cfgAPI, 5*time.Second, TIMEOUT))

	e.Logger.WithFields(map[string]interface{}{
		"numPlugins": e.numPlugin,
		"addr":       e.addr,
		"prefix":     e.prefix,
	}).Info("start listening")

	return http.ListenAndServe(e.addr, mux)
}
