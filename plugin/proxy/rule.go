// Package proxy rule.go contains proxy config struct to load and read
package proxy

type ruleType string

// const (
// 	ruleTypURI ruleType = "uri"
// 	ruleTypSRV ruleType = "srv"
// )

// Config global proxy plugin config
type Config struct {
	// reverse proxy path rules
	PathRules []PathRule `json:"path_rules"`

	// reverse proxy server rules
	ServerRules []ServerRule `json:"server_rules"`

	// ReverseServerCfgs includes all reverse servers config
	ReverseServerCfgs map[string][]ReverseServerCfg `json:"reverse_server_cfgs"`
}

// PathRule contains fields to appoint to reverse server and URL
type PathRule struct {
	// Path match input path
	Path string `json:"path"`

	// target path to reverse server
	RewritePath string `json:"rewrite_path"`

	// Method allow those Methods to request with `,` to split
	Method string `json:"method"`

	// appoint to a server to proxy only work while NeedCombine is `False`
	ServerName string `json:"server_name"`

	// CombineReqCfgs combine request settings
	CombineReqCfgs []CombineReqCfg `json:"combine_req_cfgs"`

	// need combine over two request result into one result
	NeedCombine bool `json:"need_combine"`
}

// ServerRule ...
type ServerRule struct {
	Prefix          string `json:"prefix"`
	ServerName      string `json:"server_name"`
	NeedStripPrefix bool   `json:"need_strip_prefix"`
}

// ReverseServerCfg means proxy server config.
// it contains Addr of server and Weight of current server in all servers
type ReverseServerCfg struct {
	// Name, reverse server name in config and can be used in combine config
	Name string `json:"name"`
	// Prefix, like `/rever` is usually simillar to it's Name
	Prefix string `json:"prefix"`
	// target addr
	Addr string `json:"addr"`
	// weigth for current server config
	Weight int `json:"weight"`
}

// W func return the weight of server
// implement ServerCfgInterface ...
func (s ReverseServerCfg) W() int {
	return s.Weight
}

// CombineReqCfg ...
type CombineReqCfg struct {
	ServerName string `json:"server_name"` // http://ip:port/path?params
	Path       string `json:"path"`
	Field      string `json:"field"`  // want got field
	Method     string `json:"method"` // want match method
}
