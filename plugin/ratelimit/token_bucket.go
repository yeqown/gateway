package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/yeqown/gateway/plugin"
)

var (
	_ plugin.Plugin = &Bucket{}
)

// New a Bucket to limit request with token
func New(cap, r int) *Bucket {
	if r > cap {
		panic("i panic anyway")
	}

	var (
		// once sync.Once
		b = &Bucket{
			capacity: cap,
			r:        r,
			rest:     cap / 2,
			enabled:  true,
			status:   plugin.Working,
		}
	)
	b.init()
	return b
}

// Bucket contains token for request
type Bucket struct {
	capacity int // max token to own
	r        int // speed to generate a token per second
	rest     int // rest token count
	rwm      sync.RWMutex
	enabled  bool
	status   plugin.PlgStatus
}

// init bucket
func (b *Bucket) init() {
	go b.startGenerateToken()
}

// Handle func for bucket to work as a plugin.Plugin
func (b *Bucket) Handle(c *plugin.Context) {
	if !b.accquire() {
		c.String(http.StatusTooManyRequests, "Too Many Request")
		return
	}
	c.Next()
}

// Enabled ...
func (b *Bucket) Enabled() bool {
	return b.enabled
}

// Status ...
func (b *Bucket) Status() plugin.PlgStatus {
	return b.status
}

// Name ...
func (b *Bucket) Name() string {
	return "plugin.ratelimit"
}

// Enable ...
func (b *Bucket) Enable(enabled bool) {
	b.enabled = enabled
	if !enabled {
		b.status = plugin.Stopped
	} else {
		b.status = plugin.Working
	}
}

// accquire for a token to allow request process
func (b *Bucket) accquire() bool {
	b.rwm.Lock()
	defer b.rwm.Unlock()
	if b.rest > 0 {
		b.rest--
		return true
	}
	return false
}

// generate token, always start a goroutine to process this
func (b *Bucket) startGenerateToken() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			b.rwm.Lock()
			if b.capacity >= b.rest+b.r {
				b.rest = b.capacity
			} else {
				b.rest += b.r
			}
			b.rwm.Unlock()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

}
