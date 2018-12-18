// Package main is the entry of github.com/yeqown/gateway
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yeqown/gateway/config/api"
	cpresistence "github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/cache"
	"github.com/yeqown/gateway/plugin/cache/presistence"
	"github.com/yeqown/gateway/plugin/httplog"
	"github.com/yeqown/gateway/plugin/proxy"
	"github.com/yeqown/gateway/plugin/ratelimit"
)

var (
	cfgFile = flag.String("cfgFile", "./filestore.config.json", "file store data json")
	cfgDB   = flag.String("cfgDB", "./db.config.json", "db store config json file")
)

func main() {
	flag.Parse()

	// config store
	store := cpresistence.NewJSONFileStore(*cfgFile)
	api.SetGlobal(store)
	cfg := store.Instance()

	// init logger path
	logger.InitLogger(cfg.Logpath)

	// initial plugins
	plgProxy := proxy.New(cfg.ProxyReverseServers, cfg.ProxyPathRules, cfg.ProxyServerRules)
	plgHTTPLogger := httplog.New(logger.Logger)
	plgCache := cache.New(presistence.NewInMemoryStore(), cfg.Nocache)
	plgTokenBucket := ratelimit.New(10, 1)

	e := &Engine{
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
