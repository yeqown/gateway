// Package proxy rule.go contains proxy config struct to load and read
package proxy

// Config ...
type Config struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	// Multi server configs
	Servers []Server `json:"servers"`
	// combine settings
	Combines []Combine `json:"combines"`
	// need combine over two request result into one result
	NeedCombine bool `json:"need_combine"`
}

// Server means proxy server config.
// it contains Addr of server and Weight of current server in all servers
type Server struct {
	// target addr
	Addr string `json:"addr"`
	// weigth for current server config
	Weight int `json:"weight"`
}

// W to weight of server
func (s Server) W() int {
	return s.Weight
}

// Combine ...
type Combine struct {
	Addr   string `json:"addr"`   // http://ip:port/path?params
	Field  string `json:"field"`  // want got field
	Method string `json:"method"` // want match method
}
