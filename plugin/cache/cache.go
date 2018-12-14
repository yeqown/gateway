// Package cache ... do connect to redis with RedisConfig ref to common or other where?
// declare interfaces to use cahce in common
package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/cache/presistence"
)

const (
	// CachePluginKey = "plugin.cache"
	CachePluginKey = "plugin.cache"
	// CachePageKey   = "plugin.cache.page"
	CachePageKey = "plugin.cache.page"
)

// New PluginStore ...
func New(store presistence.Store) *PluginStore {
	return &PluginStore{
		store:         store,
		serializeForm: false,
	}
}

// responseCache to save cache of one URI
// TODO: post method URI need to be cached or not? serialize the form with URI can solve this?
type responseCache struct {
	// http
	Header http.Header
	// http status code
	Status int
	// body to Save
	Data []byte
}

// responseCache encode into bytes
func encodeCache(cache *responseCache) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(cache); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// decode data byte into responseCache
func decodeToCache(byts []byte) (responseCache, error) {
	var (
		buffer bytes.Buffer
		c      responseCache
	)
	if _, err := buffer.Write(byts); err != nil {
		return c, err
	}

	dec := gob.NewDecoder(&buffer)
	if err := dec.Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}

// cachedWriter ...
type cachedWriter struct {
	http.ResponseWriter
	cache *responseCache
	store presistence.Store
	// http status code
	status int
	// key to save or get from cache
	key    string
	expire time.Duration
}

func (w cachedWriter) Header() http.Header { return w.ResponseWriter.Header() }
func (w cachedWriter) WriteHeader(status int) {
	// println("calling WriteHeader", status)
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err != nil {
		return ret, fmt.Errorf("could not write response: %v", err)
	}

	w.cache.Status = w.status
	w.cache.Header = w.Header()
	w.cache.Data = append(w.cache.Data, data...)
	// if byts, err := w.store.Get(w.key); err == nil {
	// 	w.cache.Data = append(w.cache.Data, byts...)
	// 	// fill cache other field
	// 	w.cache.Status = w.status
	// 	w.cache.Header = w.Header()
	// } else {
	// 	log.Printf("plugin.Cache Get err: %v", err)
	// }

	if w.status < 300 {
		value, err := encodeCache(w.cache)
		if err != nil {
			log.Printf("could not encode cache: %v", err)
		} else {
			if err = w.store.Set(w.key, value, w.expire); err != nil {
				log.Printf("could not set into cache: %v", err)
			}
		}
	}
	return ret, err
}

// PluginStore to serve page cache ...
type PluginStore struct {
	// store interface with value
	store presistence.Store
	// serialize form [query and post form]
	serializeForm bool
}

// TODO: support query and postform
func (p *PluginStore) generateKey(uri string) string {
	return uri
}

// func (p *PluginStore) parseValue() []byte {
// 	return []byte("hello cache")
// }

// Handle implement the interface Plugin
// [fixed] TOFIX: cannot set cache to response
func (p *PluginStore) Handle(ctx *plugin.Context) {
	defer plugin.Recover("PluginStore")

	key := p.generateKey(ctx.Path)
	// value := p.parseValue()
	if p.store.Exists(key) {
		byts, err := p.store.Get(key)
		if err != nil {
			ctx.SetError(err)
			ctx.Abort(http.StatusInternalServerError,
				fmt.Sprintf("plugin.cache Get cache err: %v", err))
			return
		}

		// decode into cache
		cache, err := decodeToCache(byts)
		if err != nil || cache.Status == 0 {
			ctx.SetError(err)
			ctx.Abort(http.StatusInternalServerError,
				fmt.Sprintf("plugin.cache decode cache err: %v", err))
			return
		}

		// write to response
		// set cache to responseWriter
		// log.Println("hit cache", string(cache.Data), cache)
		ctx.ResponseWriter().WriteHeader(cache.Status)
		ctx.ResponseWriter().Write(cache.Data)
		for k, vals := range cache.Header {
			for _, v := range vals {
				ctx.ResponseWriter().Header().Set(k, v)
			}
		}
		ctx.Abort(http.StatusOK, "")
		return
	}

	// continue process
	// println("does not hit cache")
	writer := cachedWriter{
		ResponseWriter: ctx.ResponseWriter(),
		cache:          &responseCache{},
		store:          p.store,
		status:         http.StatusOK,
		key:            key,
		expire:         presistence.DefaultExpire,
	}

	ctx.SetResponseWriter(writer)
	ctx.Next()
}
