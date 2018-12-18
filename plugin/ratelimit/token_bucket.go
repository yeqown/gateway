package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/yeqown/gateway/plugin"
)

// New a Bucket to limit request with token
func New(cap, r int) *Bucket {
	if r > cap {
		panic("i panic anyway")
	}

	var (
		// once sync.Once
		b = &Bucket{
			capacity: 0,
			r:        r,
		}
	)
	b.init()
	// once.Do(func() {
	// 	b.init()
	// })
	return b
}

// Bucket contains token for request
type Bucket struct {
	capacity int // max token to own
	r        int // speed to generate a token per second
	rest     int // rest token count
	rwm      sync.RWMutex
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
