// Package proxy ...
// this file mainly to load from file and set proxy rules
package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/utils"
)

const (
	reverseKeyLayout = "%s_%d"
	ctxTIMEOUT       = time.Second * 5
)

var (
	// ErrBalancerNotMatched plugin.Proxy balancer not matched
	ErrBalancerNotMatched = errors.New("plugin.Proxy balancer not matched")
	// ErrPageNotFound can not found page
	ErrPageNotFound = errors.New("404 Page Not Found")
)

func defaultHandleFunc(w http.ResponseWriter, req *http.Request, params httprouter.Params) {}

// New ... load configs from outter to generate a new proxy plugin
func New(c *Config) *Proxy {
	p := &Proxy{
		router:         httprouter.New(),
		balancers:      make(map[string]*Balancer),
		srvCfgsMap:     make(map[string]ReverseServerCfg),
		pathRulesMap:   make(map[string]PathRule),
		srvRulesMap:    make(map[string]ServerRule),
		reverseProxies: make(map[string]*httputil.ReverseProxy),
	}

	// initial work
	p.loadBalancers(c.ReverseServerCfgs)
	p.loadReverseProxyPathRules(c.PathRules)
	p.loadReverseProxyServerRules(c.ServerRules)

	return p
}

// Proxy ...
type Proxy struct {
	// path router
	router       *httprouter.Router
	pathRulesMap map[string]PathRule // path as key and config
	// reverse proxy configs and balancers
	reverseProxies map[string]*httputil.ReverseProxy
	srvCfgsMap     map[string]ReverseServerCfg
	balancers      map[string]*Balancer // balancer to distribute
	// server proxy configs
	srvRulesMap map[string]ServerRule
}

// Handle ... proxy to handle with request ...
// 1. single path
// 2. all server proxy
func (p *Proxy) Handle(c *plugin.Context) {
	defer plugin.Recover("Proxy")

	if p.matchedPathRule(c.Method, c.Path) {
		// callReverseURI
		logger.Logger.Info("matched path rules")
		pathRule := p.pathRulesMap[c.Path]
		if err := p.callReverseURI(pathRule, c); err != nil {
			c.SetError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if rule, ok := p.matchedServerRule(c.Path); ok {
		// callReverseServer
		logger.Logger.Info("matched server rules")
		if err := p.callReverseServer(rule, c); err != nil {
			c.SetError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// don't matched any path or server !!!
	logger.Logger.Infof("could not match path or server rule !!! (method: %s, path: %s)",
		c.Method, c.Path)
	c.SetError(ErrPageNotFound)
	c.AbortWithStatus(http.StatusNotFound)
	return
}

func (p *Proxy) matchedPathRule(method, path string) bool {
	handle, params, tsr := p.router.Lookup(method, path)
	_, _ = params, tsr
	return handle != nil
}

func (p *Proxy) matchedServerRule(path string) (ServerRule, bool) {
	pathPrefix := utils.ParseURIPrefix(path)
	pathPrefix = strings.ToLower(pathPrefix)
	rule, ok := p.srvRulesMap[pathPrefix]
	return rule, ok
}

// to load cfgs (type []proxy.ReverseServerCfg) to initial Proxy.Balancers
func (p *Proxy) loadBalancers(cfgs map[string][]ReverseServerCfg) {
	for _, cfg := range cfgs {
		srvCfgs := make([]ServerCfgInterface, len(cfg))
		for idx, srv := range cfg {
			srvCfgs[idx] = srv
			target, err := url.Parse(srv.Addr)
			if err != nil {
				panic(utils.Fstring("could not parse URL: %s", srv.Addr))
			}

			// generate reverse proxy
			reversePorxy := httputil.NewSingleHostReverseProxy(target)

			// register a func for reverse proxy to handler error
			reversePorxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
				utils.ResponseString(w, err.Error())
				return
			}
			key := utils.Fstring("%s_%d", srv.Name, idx)
			key = strings.ToLower(key)
			p.reverseProxies[key] = reversePorxy
			p.srvCfgsMap[key] = srv
		}
		// o !
		if len(cfg) != 0 {
			p.balancers[cfg[0].Name] = NewBalancer(srvCfgs)
		}
	}
}

// to load rules (type []proxy.PathRule) to initial
func (p *Proxy) loadReverseProxyPathRules(rules []PathRule) {
	for _, rule := range rules {
		// [done] TODO: valid rule all string need to be lower
		path := strings.ToLower(rule.Path)
		method := strings.ToLower(rule.Method)
		if _, ok := p.pathRulesMap[path]; ok {
			panic(utils.Fstring("duplicate path rule: %s", path))
		}
		// TODO: generate new rule with lower case string
		p.pathRulesMap[path] = rule
		for _, method := range strings.Split(rule.Method, ",") {
			p.router.Handle(method, path, defaultHandleFunc)
		}

		logger.Logger.Infof("URI rule:%s_%s registered", path, method)
	}
}

//  to load rules (type []proxy.ServerRule) to initial
func (p *Proxy) loadReverseProxyServerRules(rules []ServerRule) {
	for _, rule := range rules {
		// [done] TODO: valid rule all string need to be lower
		prefix := strings.ToLower(rule.Prefix)
		if len(prefix) <= 1 {
			log.Printf("error: prefix of %s is too short, so skipped\n", prefix)
			continue
		}
		if prefix[0] != '/' {
			prefix = "/" + prefix
		}

		if _, ok := p.srvRulesMap[prefix]; ok {
			panic(utils.Fstring("duplicate server rule prefix: %s", prefix))
		}
		// TODO: new with lower case string
		p.srvRulesMap[prefix] = rule
		logger.Logger.Infof("SRV rule:%s_%s registered", rule.ServerName, rule.Prefix)
	}
}

// callReverseURI reverse proxy to remote server and combine repsonse.
func (p *Proxy) callReverseURI(rule PathRule, c *plugin.Context) error {
	oriPath := strings.ToLower(rule.Path)
	req := c.Request()
	w := c.ResponseWriter()
	// pure reverse proxy here
	if !rule.NeedCombine {
		if len(rule.RewritePath) != 0 {
			req.URL.Path = rule.RewritePath
		}

		srvName := strings.ToLower(rule.ServerName)
		bla, ok := p.balancers[srvName]
		if !ok {
			logger.Logger.Errorf("could not found balancer of %s", oriPath)
			errmsg := utils.Fstring("error: plugin.Proxy balancer not found! (path: %s)", oriPath)
			return fmt.Errorf("%v", errmsg)
		}

		idx := bla.Distribute()
		key := strings.ToLower(
			utils.Fstring("%s_%d", srvName, idx))
		logger.Logger.Infof("proxy to server: %s URI: %s", key, req.URL.Path)

		reverseProxy, ok := p.reverseProxies[key]
		if !ok {
			return fmt.Errorf("error: plugin.Proxy reverse server not found! (key: %s)", key)
		}
		reverseProxy.ServeHTTP(w, req)
		return nil
	}
	// [done] TODO: combine two or more response
	respChan := make(chan responseChan, len(rule.CombineReqCfgs))
	ctx, cancel := context.WithTimeout(context.Background(), ctxTIMEOUT)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, cfg := range rule.CombineReqCfgs {
		wg.Add(1)
		go func(cfg CombineReqCfg) {
			defer wg.Done()
			bla, ok := p.balancers[cfg.ServerName]
			if !ok {
				respChan <- responseChan{
					Err:   ErrBalancerNotMatched,
					Field: cfg.Field,
					Data:  nil,
				}
			}
			idx := bla.Distribute()
			srvCfg, _ := p.srvCfgsMap[utils.Fstring("%s_%d", cfg.ServerName, idx)]
			combineReq(ctx, srvCfg.Addr, nil, cfg, respChan)
		}(cfg)
	}

	wg.Wait()
	close(respChan)

	r := map[string]interface{}{
		"code":    0,
		"message": "combine result",
	}

	// loop response combine to togger response
	for resp := range respChan {
		if resp.Err != nil {
			r[resp.Field] = resp.Err.Error()
			continue
		}
		// read response
		r[resp.Field] = resp.Data
	}

	// Response
	c.JSON(http.StatusOK, r)

	return nil
}

// callReverseServer to proxy request to another server
// cannot combine two server response
func (p *Proxy) callReverseServer(rule ServerRule, c *plugin.Context) error {
	// need to trim prefix
	req := c.Request()
	w := c.ResponseWriter()
	if rule.NeedStripPrefix {
		req.URL.Path = strings.TrimPrefix(strings.ToLower(req.URL.Path),
			strings.ToLower(rule.Prefix))
	}

	srvName := strings.ToLower(rule.ServerName)
	bla, ok := p.balancers[srvName]

	if !ok {
		errmsg := utils.Fstring("%s Not Found!", srvName)
		return fmt.Errorf("%v", errmsg)
	}

	idx := bla.Distribute()
	key := utils.Fstring("%s_%d", srvName, idx)
	key = strings.ToLower(key)

	reverseProxy, ok := p.reverseProxies[key]
	if !ok {
		return fmt.Errorf("could not found reverse proxy")
	}
	logger.Logger.Infof("proxy to %s", req.URL.Path)
	reverseProxy.ServeHTTP(w, req)
	return nil
}
