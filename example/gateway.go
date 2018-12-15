// Package main is the entry of github.com/yeqown/gateway
package main

import (
	"log"

	"github.com/yeqown/gateway"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/cache"
	"github.com/yeqown/gateway/plugin/cache/presistence"
	"github.com/yeqown/gateway/plugin/httplog"
	"github.com/yeqown/gateway/plugin/proxy"
)

var (
	c            proxy.Config
	nocacheRules []cache.Rule
)

func init() {
	c = proxy.Config{
		PathRules: []proxy.PathRule{
			proxy.PathRule{
				Path:           "/gw/name",
				RewritePath:    "/srv/name",
				Method:         "GET",
				ServerName:     "srv1",
				CombineReqCfgs: []proxy.CombineReqCfg{}, // no need if NeedCombine false
				NeedCombine:    false,
			},
			proxy.PathRule{
				Path:           "/gw/id",
				RewritePath:    "/srv/id",
				Method:         "POST,GET",
				ServerName:     "srv1",
				CombineReqCfgs: []proxy.CombineReqCfg{},
				NeedCombine:    false,
			},
			proxy.PathRule{
				Path:        "/gw/combine",
				RewritePath: "",    // no need
				Method:      "GET", // combine Method
				ServerName:  "",    // no need
				CombineReqCfgs: []proxy.CombineReqCfg{
					proxy.CombineReqCfg{
						ServerName: "srv1",
						Path:       "/srv/id",
						Field:      "id",
						Method:     "GET",
					},
					proxy.CombineReqCfg{
						ServerName: "srv1",
						Path:       "/srv/name",
						Field:      "name",
						Method:     "POST",
					},
				},
				NeedCombine: true,
			},
		},
		ServerRules: []proxy.ServerRule{
			proxy.ServerRule{
				Prefix:          "/srv",
				ServerName:      "srv1",
				NeedStripPrefix: false,
			},
			proxy.ServerRule{
				Prefix:          "/srvPrefix",
				ServerName:      "srv1",
				NeedStripPrefix: true,
			},
		},
		ReverseServerCfgs: map[string][]proxy.ReverseServerCfg{
			"group1": []proxy.ReverseServerCfg{
				proxy.ReverseServerCfg{
					Name:   "srv1",
					Prefix: "/srv",
					Addr:   "http://127.0.0.1:8081",
					Weight: 5,
				},
				proxy.ReverseServerCfg{
					Name:   "srv1",
					Prefix: "/srv",
					Addr:   "http://127.0.0.1:8082",
					Weight: 5,
				},
			},
		},
	}

	nocacheRules = []cache.Rule{
		cache.Rule{
			Regular: "^/gw/id$",
		},
	}
}

// main will do some initalize work and
// start the whole api-gateway
func main() {
	logger.InitLogger("./logs")

	// initial plugins
	plgProxy := proxy.New(&c)
	plgHTTPLogger := httplog.New(logger.Logger)
	plgCache := cache.New(presistence.NewInMemoryStore(), nocacheRules)

	eng := &gateway.Engine{
		Logger: logger.Logger,
		Plugins: []plugin.Plugin{
			plgHTTPLogger,
			plgCache,
			plgProxy,
		},
	}

	if err := eng.ListenAndServe(":8989"); err != nil {
		log.Fatal(err)
	}
}
