package gateway

import (
	"net/http"
	"strings"
	"time"

	"github.com/yeqown/gateway/plugin"
	log "github.com/yeqown/server-common/logger"
)

// Engine ...
type Engine struct {
	Plugins   []plugin.Plugin
	numPlugin int
	addr      string

	Prefix string
	Logger *log.Logger // inner logger
}

func (e *Engine) init() {
	e.numPlugin = len(e.Plugins)

	if len(e.Prefix) <= 1 {
		e.Prefix = "/api/"
	}

	if e.Prefix[0] != '/' {
		e.Prefix = "/" + e.Prefix
	}

	if e.Prefix[len(e.Prefix)-1] != '/' {
		e.Prefix = e.Prefix + "/"
	}
}

func (e *Engine) use(plgs ...plugin.Plugin) {
	e.Plugins = append(e.Plugins, plgs...)
	e.numPlugin += len(plgs)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path,
		strings.TrimSuffix(e.Prefix, "/"))
	// ctx := ctxPool.Get().(*plugin.Context)
	// defer ctxPool.Put(ctx)
	e.Logger.Info("new request recved")

	// generate a new context
	ctx := plugin.New(w, req, e.numPlugin, e.Plugins)

	// start call
	ctx.Next()

	// reset resource
	ctx.Reset()
	return
}

// ListenAndServe ...
func (e *Engine) ListenAndServe(addr string) error {
	e.addr = addr
	e.init()

	mux := http.NewServeMux()

	mux.Handle(webPrefix, HTMLSrv)
	mux.Handle(e.Prefix, http.TimeoutHandler(e, 5*time.Second, "timeout"))

	e.Logger.WithFields(map[string]interface{}{
		"numPlugins": e.numPlugin,
		"addr":       e.addr,
		"prefix":     e.Prefix,
	}).Info("start listening")

	return http.ListenAndServe(addr, mux)
}
