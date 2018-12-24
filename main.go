// Package main is the entry of github.com/yeqown/gateway
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yeqown/gateway/config/api"
	"github.com/yeqown/gateway/config/presistence/mongostore"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin/cache"
	"github.com/yeqown/gateway/plugin/cache/presistence"
	"github.com/yeqown/gateway/plugin/httplog"
	"github.com/yeqown/gateway/plugin/proxy"
	"github.com/yeqown/gateway/plugin/ratelimit"
)

var (
	logpath = flag.String("logpath", "./logs", "whcih path will log files be stored ")
	port    = flag.Int("port", 8989, "the config gate server will listen at the port")
	// cfgDB   = flag.String("cfgDB", "./db.config.json", "db store config json file")
)

func main() {
	flag.Parse()

	// config store
	store, err := mongostore.New("mongodb://127.0.0.1:27017", "gateway")
	if err != nil {
		panic(err)
	}

	api.SetGlobal(store)
	cfg := store.Instance()

	// init logger path
	logger.InitLogger(*logpath)

	// initial plugins
	plgProxy := proxy.New(cfg.ProxyReverseServers, cfg.ProxyPathRules, cfg.ProxyServerRules)
	plgHTTPLogger := httplog.New(logger.Logger)
	plgCache := cache.New(presistence.NewInMemoryStore(), cfg.Nocache)
	plgTokenBucket := ratelimit.New(10, 1)

	e := &Engine{
		Logger:        logger.Logger,
		StoreChangedC: store.Updated(),
	}

	// TODO: load with active plugin names list
	e.use(plgHTTPLogger, plgTokenBucket, plgCache, plgProxy)

	addr := fmt.Sprintf(":%d", *port)
	if err := e.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
