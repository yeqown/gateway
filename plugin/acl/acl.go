// Package acl means "Access Control List"
// to Limit traffic and some IP maybe discern spider ?
package acl

import "time"

// RequestAccessLimiter ...
type RequestAccessLimiter struct {
	ipWhiteList       map[string][]string
	RequestTimesLimit map[string]*timesLimit
}

type timesLimit struct {
	ip              string
	max             int
	cnt             int
	lastRequestTime time.Time
}
