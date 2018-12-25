// Package cache ... do connect to redis with RedisConfig ref to common or other where?
// declare interfaces to use cahce in common
package cache

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/logger"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/plugin/cache/presistence"
	"github.com/yeqown/gateway/utils"
)

var (
	_ plugin.Plugin = &Cache{}
)

const (
	// CachePluginKey = "plugin.cache"
	CachePluginKey = "plugin.cache"
	// CachePageKey   = "plugin.cache.page"
	CachePageKey = "plugin.cache.page"
)

// New PluginStore ...
func New(store presistence.Store, rules []rule.Nocacher) *Cache {
	c := &Cache{
		store:         store,
		serializeForm: false,
		status:        plugin.Working,
		enabled:       true,
	}
	c.Load(rules)
	return c
}

// responseCache to save cache of one URI
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
	cache  *responseCache
	store  presistence.Store // http status code
	status int               // key to save or get from cache
	key    string
	expire time.Duration
}

func (w cachedWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w cachedWriter) WriteHeader(status int) {
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

	// only save into store while statusCode is successful
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

// Cache to serve page cache ...
type Cache struct {
	store         presistence.Store // store interface with value
	serializeForm bool              // serialize form [query and post form]
	enabled       bool
	status        plugin.PlgStatus
	regexps       []*regexp.Regexp // store regular expression
	cntRegexp     int              // count of regexps
}

// generate a key with the given http.Request and serializeForm flag
// [done] TODO: post method URI need to be cached or not? serialize the form with URI can solve this?
func generateKey(req *http.Request, serializeForm bool) string {
	var (
		cpyReq     *http.Request
		formEncode string
	)
	if serializeForm {
		cpyReq = utils.CopyRequest(req)
		formEncode = utils.EncodeFormToString(cpyReq)
		return urlEscape(CachePluginKey, req.URL.RequestURI(), formEncode)
	}

	return urlEscape(CachePluginKey, req.URL.RequestURI())
}

func urlEscape(prefix, u string, extern ...string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	for _, s := range extern {
		buffer.WriteString(":")
		buffer.WriteString(s)
	}
	return buffer.String()
}

// Handle implement the interface Plugin
// [fixed] TOFIX: cannot set cache to response
func (c *Cache) Handle(ctx *plugin.Context) {
	defer plugin.Recover("Cache")

	if c.matchNoCacheRule(ctx.Path) {
		logger.Logger.Infof("plugin.Cache cannot work with path: %s", ctx.Path)
		return
	}

	logger.Logger.Info("plugin.Cache is working")
	key := generateKey(ctx.Request(), c.serializeForm)
	if c.store.Exists(key) {
		// if exists key then load from cache and then
		// write to http.ResponseWriter
		byts, err := c.store.Get(key)
		if err != nil {
			ctx.SetError(fmt.Errorf("plugin.cache Get cache err: %v", err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// decode into cache
		cache, err := decodeToCache(byts)
		if err != nil || cache.Status == 0 {
			ctx.SetError(fmt.Errorf("plugin.cache decode cache err: %v", err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
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
		ctx.AbortWithStatus(http.StatusOK)
		return
	}

	// continue process
	// println("does not hit cache")
	writer := cachedWriter{
		ResponseWriter: ctx.ResponseWriter(),
		cache:          &responseCache{},
		store:          c.store,
		status:         http.StatusOK,
		key:            key,
		expire:         presistence.DefaultExpire,
	}

	ctx.SetResponseWriter(writer)
	ctx.Next()
}

// Enabled ...
func (c *Cache) Enabled() bool {
	return c.enabled
}

// Status ...
func (c *Cache) Status() plugin.PlgStatus {
	return c.status
}

// Name ...
func (c *Cache) Name() string {
	return "plugin.cache"
}

// Enable ...
func (c *Cache) Enable(enabled bool) {
	c.enabled = enabled
	if !enabled {
		c.status = plugin.Stopped
	} else {
		c.status = plugin.Working
	}
}
