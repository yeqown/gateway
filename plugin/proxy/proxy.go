// Package proxy ...
// this file mainly to load from file and set proxy rules
package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/utils"
)

func hrHandleF(w http.ResponseWriter, req *http.Request, params httprouter.Params) {}

// New ... load configs from outter to generate a new proxy plugin
func New(proxyCfgs []Config) *Proxy {
	p := &Proxy{
		cfgs:           proxyCfgs,
		router:         httprouter.New(),
		balancers:      make(map[string]*Balancer),
		pathCfgs:       make(map[string]Config),
		reverseProxies: make(map[string]*httputil.ReverseProxy),
	}
	p.load()
	return p
}

// Proxy ...
type Proxy struct {
	cfgs           []Config
	router         *httprouter.Router
	balancers      map[string]*Balancer // balancer to distribute
	pathCfgs       map[string]Config    // path as key and config
	reverseProxies map[string]*httputil.ReverseProxy
}

// Handle ... proxy to handle with request ...
func (p *Proxy) Handle(c *plugin.Context) {
	defer plugin.Recover("Proxy")

	handle, params, tsr := p.router.Lookup(c.Method, c.Path)
	_, _ = params, tsr

	if handle == nil {
		msg := utils.Fstring("%s %s error! tip: Not Found", c.Method, c.Path)
		c.Abort(http.StatusNotFound, msg)
		return
	}

	// Call server
	if err := p.Call(c.Path, c.ResponseWriter(), c.Request()); err != nil {
		c.Abort(http.StatusInternalServerError, err.Error())
		return
	}
}

func (p *Proxy) load() {
	for _, cfg := range p.cfgs {
		// init balancers ...
		srvCfgs := make([]ServerCfg, len(cfg.Servers))
		for idx, srv := range cfg.Servers {
			srvCfgs[idx] = srv
			target, err := url.Parse(srv.Addr)
			if err != nil {
				panic(err)
			}
			prefix := utils.Fstring("%s_%d", cfg.Path, idx)

			// generate reverse proxy
			reversePorxy := httputil.NewSingleHostReverseProxy(target)

			// handler error
			reversePorxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
				utils.ResponseString(w, err.Error())
				return
			}
			p.reverseProxies[prefix] = reversePorxy
		}
		// init balancer
		p.balancers[cfg.Path] = NewBalancer(srvCfgs)
		p.pathCfgs[cfg.Path] = cfg

		// reg into router
		p.router.Handle(cfg.Method, cfg.Path, hrHandleF)
		logger.Logger.Infof("URI %s_%s registered", cfg.Path, cfg.Method)
	}
}

// Call server and repsonse
func (p *Proxy) Call(path string, w http.ResponseWriter, req *http.Request) error {
	bla, ok := p.balancers[path]
	if !ok {
		errmsg := utils.Fstring("%s Not Found!", path)
		return fmt.Errorf("%v", errmsg)

	}
	cfg := p.pathCfgs[path]
	idx := bla.Distribute()
	// prefix := utils.Fstring("%s_%d", path, idx)
	// srv := p.reverseProxies[prefix]

	// pure reverse proxy here
	if !cfg.NeedCombine {
		prefix := utils.Fstring("%s_%d", cfg.Path, idx)
		logger.Logger.Infof("proxy to %s", prefix)
		reverseProxy, ok := p.reverseProxies[prefix]
		if !ok {
			return fmt.Errorf("could not found reverse proxy")
		}
		reverseProxy.ServeHTTP(w, req)
		return nil
	}

	return nil
}
