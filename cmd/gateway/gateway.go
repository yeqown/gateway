// Package main is the entry of github.com/yeqown/gateway
package main

import (
	"log"
	"net/http"

	"github.com/yeqown/gateway"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/httplog"
	"github.com/yeqown/gateway/plugin/proxy"
)

// main will do some initalize work and
// start the whole api-gateway
func main() {
	logger.InitLogger("./logs")

	plgProxy := proxy.New([]proxy.Config{
		proxy.Config{
			Path:   "/gw/srv1",
			Method: http.MethodGet,
			Servers: []proxy.Server{
				{
					Addr:   "http://localhost:8081/srv/name",
					Weight: 2,
				},
			},
			Combines:    nil,
			NeedCombine: false,
		},
		proxy.Config{
			Path:   "/gw/srv2",
			Method: http.MethodGet,
			Servers: []proxy.Server{
				{
					Addr:   "http://localhost:8082/srv/name",
					Weight: 4,
				},
			},
			Combines:    nil,
			NeedCombine: false,
		},
	})

	plgHTTPLogger := httplog.New(logger.Logger)

	eng := &gateway.Engine{
		Logger: logger.Logger,
		Plugins: []plugin.Plugin{
			plgHTTPLogger,
			plgProxy,
		},
		Prefix: "/proxy/",
	}

	if err := eng.ListenAndServe(":8989"); err != nil {
		log.Fatal(err)
	}
}
