package rule

// ProxyManager ...
type ProxyManager interface{}

// PathRuler 用于单个url配置
type PathRuler interface {
	Path() string
	Method() string
	ServerName() string
	RewritePath() string
	// 组合请求
	NeedCombine() bool
	// 组合请求配置
	CombineReqCfgs() []Combiner
}

// ServerRuler 用于配置一组服务代理
type ServerRuler interface {
	Prefix() string
	ServerName() string
	NeedStripPrefix() bool
}

// ReverseServer 单向代理服务器配置
type ReverseServer interface {
	Name() string
	Addr() string
	W() int
}

// Combiner 用于合并请求时候用的配置
type Combiner interface {
	ServerName() string
	Path() string
	Field() string
	Method() string
}

// PathCfg contains fields to appoint to reverse server and URL
type PathCfg struct {
	PPath        string           `json:"path"`
	RRewritePath string           `json:"rewrite_path"`
	MMethod      string           `json:"method"`
	SrvName      string           `json:"server_name"`
	CombReqs     []*CombineReqCfg `json:"combine_req_cfgs"`
	NeedComb     bool             `json:"need_combine"`
}

// Path of PathCfg
func (c *PathCfg) Path() string { return c.PPath }

// Method of PathCfg
func (c *PathCfg) Method() string { return c.MMethod }

// ServerName of PathCfg
func (c *PathCfg) ServerName() string { return c.SrvName }

// RewritePath of PathCfg
func (c *PathCfg) RewritePath() string { return c.RRewritePath }

// NeedCombine of PathCfg
func (c *PathCfg) NeedCombine() bool { return c.NeedComb }

// CombineReqCfgs of PathCfg
func (c *PathCfg) CombineReqCfgs() []Combiner {
	combs := make([]Combiner, len(c.CombReqs))
	for idx, c := range c.CombReqs {
		combs[idx] = c
	}
	return combs
}

// ServerCfg ... ServerCfg ServerCfg ServerCfg ServerCfg
type ServerCfg struct {
	PPrefix          string `json:"prefix"`
	SServerName      string `json:"server_name"`
	NNeedStripPrefix bool   `json:"need_strip_prefix"`
}

// Prefix ...
func (s *ServerCfg) Prefix() string { return s.PPrefix }

// ServerName ...
func (s *ServerCfg) ServerName() string { return s.SServerName }

// NeedStripPrefix ...
func (s *ServerCfg) NeedStripPrefix() bool { return s.NNeedStripPrefix }

// ReverseServerCfg means proxy server config.
// it contains Addr of server and Weight of current server in all servers
type ReverseServerCfg struct {
	NName string `json:"name"`
	// PPrefix string `json:"prefix"`
	AAddr  string `json:"addr"`
	Weight int    `json:"weight"`
}

// Name ...
func (s ReverseServerCfg) Name() string { return s.NName }

// Addr ...
func (s ReverseServerCfg) Addr() string { return s.AAddr }

// W func return the weight of server implement ServerCfgInterface ...
func (s ReverseServerCfg) W() int { return s.Weight }

// CombineReqCfg ... CombineReqCfg CombineReqCfg CombineReqCfg CombineReqCfg
type CombineReqCfg struct {
	SrvName string `json:"server_name"` // http://ip:port/path?params
	PPath   string `json:"path"`        // path `/request/path`
	FField  string `json:"field"`       // want got field
	MMethod string `json:"method"`      // want match method
}

// ServerName () string
func (c *CombineReqCfg) ServerName() string { return c.SrvName }

// Path ...
func (c *CombineReqCfg) Path() string { return c.PPath }

// Field ...
func (c *CombineReqCfg) Field() string { return c.FField }

// Method ...
func (c *CombineReqCfg) Method() string { return c.MMethod }
