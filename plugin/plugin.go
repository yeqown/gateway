package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	// handle err happend
	if c.err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Errorf("could not handle with request, err: %v", c.err).Error())
		return
	}

	// handle aborrted
	if c.aborted {
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

// Abort process to stop calling next plugin
// [done] TODO: ignore response here, should call JSON, or String manually
func (c *Context) Abort() {
	c.aborted = true
}

// AbortWithStatus abort process and set response status
func (c *Context) AbortWithStatus(status int) {
	c.aborted = true
	c.w.WriteHeader(status)
}

// Aborted ...
func (c *Context) Aborted() bool {
	return c.aborted
}

// Set set request and  responseWriter
func (c *Context) Set(req *http.Request, w http.ResponseWriter) {
	c.req = req
	c.w = w
	c.pluginIdx = -1
}

// Reset ... donot call this manually
func (c *Context) Reset() {
	c.req = nil
	c.w = nil
	c.pluginIdx = -1
}

// Error get the global error of context
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

// JSON ...
func (c *Context) JSON(status int, v interface{}) {
	byts, err := json.Marshal(v)
	if err != nil {
		c.SetError(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(c.w, string(byts))
	c.AbortWithStatus(status)
	// c.w.WriteHeader(status)
	// c.Abort()
}

// String ...
func (c *Context) String(status int, s string) {
	fmt.Fprintf(c.w, s)
	// c.w.WriteHeader(status)
	// c.Abort()
	c.AbortWithStatus(status)
}
