package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yeqown/gateway/utils"
)

// Plugin type Plugin want to save all plugin
type Plugin interface {
	Handle(ctx *Context)
}

// New generate a Context
// TODO: do this work with pool
func New(w http.ResponseWriter, req *http.Request,
	numPlugin int, plugins []Plugin,
) *Context {
	method := req.Method
	path := req.URL.Path

	return &Context{
		Method:    method,
		Path:      path,
		numPlugin: numPlugin,
		plugins:   plugins,
		pluginIdx: -1,
		w:         w,
		req:       req,
		aborted:   false,
	}
}

// Context ... contains infomation to transfer
type Context struct {
	Ctx    context.Context // ctx control signal for multi goroutine
	Method string          // request method
	Path   string          // request Path

	req *http.Request
	w   http.ResponseWriter

	plugins   []Plugin
	pluginIdx int
	numPlugin int

	aborted bool  // request aborted
	err     error // error
}

// Next ...
func (c *Context) Next() {
	// handle aborrted
	if c.aborted {
		return
	}

	// handle err happend
	if c.err != nil {
		c.Abort(http.StatusInternalServerError,
			fmt.Errorf("could not handle with request, err: %v", c.err).Error())
		return
	}

	// call next
	c.pluginIdx++
	if c.pluginIdx >= c.numPlugin {
		return
	}

	c.plugins[c.pluginIdx].Handle(c)
	c.Next()
}

// Abort ...
func (c *Context) Abort(status int, msg string) {
	c.aborted = true
	c.w.WriteHeader(status)

	// json
	if json.Valid([]byte(msg)) {
		utils.ResponseJSON(c.w, msg)
		return
	}

	// normal string
	utils.ResponseString(c.w, msg)
}

// Aborted ...
func (c *Context) Aborted() bool {
	return c.aborted
}

// Set ...
func (c *Context) Set(req *http.Request, w http.ResponseWriter) {
	c.req = req
	c.w = w
	c.pluginIdx = -1
}

// Reset ...
func (c *Context) Reset() {
	c.req = nil
	c.w = nil
	c.pluginIdx = -1
}

// Error
func (c *Context) Error() error {
	return c.err
}

// SetError ...
func (c *Context) SetError(err error) {
	c.err = err
}

// Request ...
func (c *Context) Request() *http.Request {
	return c.req
}

// ResponseWriter ...
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

// SetResponseWriter ...
func (c *Context) SetResponseWriter(w http.ResponseWriter) {
	c.w = w
}
