// Package cache ... do connect to redis with RedisConfig ref to common or other where?
// // declare interfaces to use cahce in common
package cache

import (
	"net/http"
	"time"
)

const (
	// CachePluginKey = "plugin.cache"
	CachePluginKey = "plugin.cache"
	// CachePageKey   = "plugin.cache.page"
	CachePageKey = "plugin.cache.page"
	// DefaultExpire  = 5 * time.Minute
	DefaultExpire = 5 * time.Minute
)

type responseCache struct {
	Header http.Header
	Status int
	Data   []byte
}

// ResponseCacheWriter ...
type ResponseCacheWriter struct {
	Status int
}
