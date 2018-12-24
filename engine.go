package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config"
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/plugin"
	log "github.com/yeqown/server-common/logger"
)

// TIMEOUT string
const TIMEOUT = "timeout"

// Engine ...
type Engine struct {
	allPlugins   []plugin.Plugin // all register plugins
	numAllPlugin int             // num of plugin
	// activePlugins     []plugin.Plugin // plugins in use
	// numActivePlugin   int             // num of active plugin
	// activePluginsName []string        // 已经被启用的插件序列，用于持久化存储和加载

	addr   string      // gate addr
	prefix string      // Engine gate prefix to serve
	Logger *log.Logger // inner logger

	plgapi *httprouter.Router // plugin manage api router

	// cfgAPI        *config.HTTP                   // config api handler
	StoreChangedC <-chan presistence.ChangedChan // store changed channel
}

func (e *Engine) use(plgs ...plugin.Plugin) {
	e.allPlugins = append(e.allPlugins, plgs...)
	e.numAllPlugin += len(plgs)
}

// 启用或者禁用插件
//
func (e *Engine) enablePlugin(enabled bool, pluginIdx int) {
	// Proxy 插件不允许禁用
	if pluginIdx == (e.numAllPlugin - 1) {
		return
	}
	e.allPlugins[pluginIdx].Enable(enabled)
}

func (e *Engine) init(addr string) {
	if addr != "" {
		e.addr = addr
	}
	e.numAllPlugin = len(e.allPlugins)
	e.prefix = "/gate"

	// e.initActivePlugins()
	e.initPluginManageRouter()

	go func() {
		for {
			select {
			case c := <-e.StoreChangedC:
				e.Logger.Infof("store changed: %v", c)
				switch c.Code {
				case presistence.PlgCodeCache:
					e.Logger.Info("reload cache rules")
				case presistence.PlgCodeProxyPath:
					e.Logger.Info("reload Proxy path rules")
				case presistence.PlgCodeProxyServer:
					e.Logger.Info("reload Proxy server rulese")
				case presistence.PlgCodeProxyReverseSrv:
					e.Logger.Info("reload Proxy reversesrv rules")
				case presistence.PlgCodeRatelimit:
					e.Logger.Info("reload ratelimit rules")
				}
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
}

// func (e *Engine) initActivePlugins() {
// 	e.activePlugins = make([]plugin.Plugin, 0, e.numAllPlugin)
// 	e.activePluginsName = make([]string, 0, e.numAllPlugin)
// 	e.numActivePlugin = 0
// 	for _, plg := range e.allPlugins {
// 		// e.pluginStatus[idx] = plg.Status()
// 		if plg.Enabled() {
// 			e.numActivePlugin++
// 			e.activePlugins = append(e.activePlugins, plg)
// 			e.activePluginsName = append(e.activePluginsName, plg.Name())
// 		}
// 	}
// }

func (e *Engine) initPluginManageRouter() {
	e.plgapi = httprouter.New()
	e.plgapi.GET("/plugins", e.pluginsGET)
	e.plgapi.GET("/plugins/enable", e.enablePluginGET)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path, e.prefix)
	handle, params, tsr := e.plgapi.Lookup(req.Method, req.URL.Path)
	_ = tsr
	if handle != nil {
		handle(w, req, params)
		return
	}

	// ctx := ctxPool.Get().(*plugin.Context)
	// defer ctxPool.Put(ctx)
	ctx := plugin.New(w, req, e.numAllPlugin, e.allPlugins)
	ctx.Next()

	return
}

// ListenAndServe ...
func (e *Engine) ListenAndServe(addr string) error {
	e.init(addr)

	// init plugins mux
	mux := http.NewServeMux()
	mux.Handle("/gate/",
		http.TimeoutHandler(e, 5*time.Second, TIMEOUT))

	// init config api mux
	cfgAPI := config.New("/gateapi")
	mux.Handle("/gateapi/",
		http.TimeoutHandler(cfgAPI, 5*time.Second, TIMEOUT))

	e.Logger.WithFields(map[string]interface{}{
		"numPlugins": e.numAllPlugin,
		"addr":       e.addr,
		"prefix":     "/gate/",
	}).Info("start listening")

	// return gracehttp.Serve(
	// 	&http.Server{Addr: e.addr, Handler: mux},
	// )
	return http.ListenAndServe(e.addr, mux)
}
