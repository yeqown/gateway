// Package main is the entry of github.com/yeqown/gateway
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yeqown/gateway"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/cache"
	"github.com/yeqown/gateway/plugin/cache/presistence"
	"github.com/yeqown/gateway/plugin/httplog"
	"github.com/yeqown/gateway/plugin/proxy"
	"github.com/yeqown/gateway/plugin/ratelimit"
)

var (
	cfgFile = flag.String("cfgFile", "./config.json", "load config from this file")
)

// main will do some initalize work and
// start the whole api-gateway
func main() {
	// load config form file
	cfg, err := gateway.LoadConfig(*cfgFile)
	if err != nil {
		panic(err)
	}

	// init logger path
	logger.InitLogger(cfg.Logpath)

	// initial plugins
	plgProxy := proxy.New(cfg.ProxyCfg)
	plgHTTPLogger := httplog.New(logger.Logger)
	plgCache := cache.New(presistence.NewInMemoryStore(), cfg.CachenoRules)
	plgTokenBucket := ratelimit.New(10, 1)

	e := &gateway.Engine{
		Logger: logger.Logger,
		Plugins: []plugin.Plugin{
			plgHTTPLogger,
			plgTokenBucket,
			plgCache,
			plgProxy,
		},
	}
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := e.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
