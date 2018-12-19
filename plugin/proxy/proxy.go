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
	"github.com/yeqown/gateway/config/rule"
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
	// ErrNoReverseServer ...
	ErrNoReverseServer = errors.New("could not found reverse proxy")
)

func defaultHandleFunc(w http.ResponseWriter, req *http.Request, params httprouter.Params) {}

// New ... load configs from outter to generate a new proxy plugin
func New(
	reverseServers map[string][]rule.ReverseServer,
	pathRules []rule.PathRuler,
	srvRules []rule.ServerRuler,
) *Proxy {
	p := &Proxy{
		router:         httprouter.New(),
		balancers:      make(map[string]*Balancer),
		srvCfgsMap:     make(map[string]rule.ReverseServer),
		pathRulesMap:   make(map[string]rule.PathRuler),
		srvRulesMap:    make(map[string]rule.ServerRuler),
		reverseProxies: make(map[string]*httputil.ReverseProxy),
	}

	// initial work
	p.loadBalancers(reverseServers)
	p.loadReverseProxyPathRules(pathRules)
	p.loadReverseProxyServerRules(srvRules)

	return p
}

// Proxy ...
type Proxy struct {
	// router
	router *httprouter.Router

	// reverse proxy configs and balancers
	reverseProxies map[string]*httputil.ReverseProxy
	balancers      map[string]*Balancer // balancer to distribute

	// config
	pathRulesMap map[string]rule.PathRuler // path as key and config
	srvRulesMap  map[string]rule.ServerRuler
	srvCfgsMap   map[string]rule.ReverseServer
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

func (p *Proxy) matchedServerRule(path string) (rule.ServerRuler, bool) {
	pathPrefix := utils.ParseURIPrefix(path)
	pathPrefix = strings.ToLower(pathPrefix)
	rule, ok := p.srvRulesMap[pathPrefix]
	return rule, ok
}

// to load cfgs (type []proxy.ReverseServerCfg) to initial Proxy.Balancers
func (p *Proxy) loadBalancers(cfgs map[string][]rule.ReverseServer) {
	for _, cfg := range cfgs {
		srvCfgs := make([]ServerCfgInterface, len(cfg))
		for idx, srv := range cfg {
			srvCfgs[idx] = srv
			target, err := url.Parse(srv.Addr())
			if err != nil {
				panic(utils.Fstring("could not parse URL: %s", srv.Addr()))
			}

			// generate reverse proxy
			reversePorxy := httputil.NewSingleHostReverseProxy(target)

			// register a func for reverse proxy to handler error
			reversePorxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
				utils.ResponseString(w, err.Error())
				return
			}
			key := utils.Fstring("%s_%d", srv.Name(), idx)
			key = strings.ToLower(key)
			p.reverseProxies[key] = reversePorxy
			p.srvCfgsMap[key] = srv
		}
		// o !
		if len(cfg) != 0 {
			p.balancers[cfg[0].Name()] = NewBalancer(srvCfgs)
		}
	}
}

// to load rules (type []proxy.PathRule) to initial
func (p *Proxy) loadReverseProxyPathRules(rules []rule.PathRuler) {
	for _, rule := range rules {
		// [done] TODO: valid rule all string need to be lower
		path := strings.ToLower(rule.Path())
		method := strings.ToLower(rule.Method())
		if _, ok := p.pathRulesMap[path]; ok {
			panic(utils.Fstring("duplicate path rule: %s", path))
		}
		p.pathRulesMap[path] = rule
		for _, method := range strings.Split(rule.Method(), ",") {
			p.router.Handle(method, path, defaultHandleFunc)
		}

		logger.Logger.Infof("URI rule:%s_%s registered", path, method)
	}
}

//  to load rules (type []proxy.ServerRule) to initial
func (p *Proxy) loadReverseProxyServerRules(rules []rule.ServerRuler) {
	for _, rule := range rules {
		// [done] TODO: valid rule all string need to be lower
		prefix := strings.ToLower(rule.Prefix())
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
		p.srvRulesMap[prefix] = rule
		logger.Logger.Infof("SRV rule:%s_%s registered", rule.ServerName(), rule.Prefix())
	}
}

// callReverseURI reverse proxy to remote server and combine repsonse.
func (p *Proxy) callReverseURI(pr rule.PathRuler, c *plugin.Context) error {
	oriPath := strings.ToLower(pr.Path())
	req := c.Request()
	w := c.ResponseWriter()
	// pure reverse proxy here
	if !pr.NeedCombine() {
		if len(pr.RewritePath()) != 0 {
			req.URL.Path = pr.RewritePath()
		}

		srvName := strings.ToLower(pr.ServerName())
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
	respChan := make(chan responseChan, len(pr.CombineReqCfgs()))
	ctx, cancel := context.WithTimeout(context.Background(), ctxTIMEOUT)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, combCfg := range pr.CombineReqCfgs() {
		wg.Add(1)
		go func(comb rule.Combiner) {
			defer wg.Done()
			bla, ok := p.balancers[comb.ServerName()]
			if !ok {
				respChan <- responseChan{
					Err:   ErrBalancerNotMatched,
					Field: comb.Field(),
					Data:  nil,
				}
			}
			idx := bla.Distribute()
			srvCfg, _ := p.srvCfgsMap[utils.Fstring("%s_%d", comb.ServerName(), idx)]
			combineReq(ctx, srvCfg.Addr(), nil, comb, respChan)
		}(combCfg)
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
func (p *Proxy) callReverseServer(sr rule.ServerRuler, c *plugin.Context) error {
	// need to trim prefix
	req := c.Request()
	w := c.ResponseWriter()
	if sr.NeedStripPrefix() {
		req.URL.Path = strings.TrimPrefix(strings.ToLower(req.URL.Path),
			strings.ToLower(sr.Prefix()))
	}

	srvName := strings.ToLower(sr.ServerName())
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
		return ErrNoReverseServer
	}
	logger.Logger.Infof("proxy to %s", req.URL.Path)
	reverseProxy.ServeHTTP(w, req)
	return nil
}
