package filestore

// filerule.go are struct implement rule/*.Interface
// and used by filestore.go

import (
	"fmt"

	"github.com/yeqown/gateway/config/rule"
)

// PathCfg contains fields to appoint to reverse server and URL
type PathCfg struct {
	PPath        string           `json:"path"`
	RRewritePath string           `json:"rewrite_path"`
	MMethod      string           `json:"method"`
	SrvName      string           `json:"server_name"`
	CombReqs     []*CombineReqCfg `json:"combine_req_cfgs"`
	NeedComb     bool             `json:"need_combine"`
	Idx          string           `json:"-"`
}

// ID func implement Ruler interface
func (c *PathCfg) ID() string {
	return c.Idx
}

// String func implement Ruler interface
func (c *PathCfg) String() string {
	return fmt.Sprintf("id: %s, path: %s", c.Idx, c.PPath)
}

// SetID func implement Ruler interface
func (c *PathCfg) SetID(id string) {
	c.Idx = id
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
func (c *PathCfg) CombineReqCfgs() []rule.Combiner {
	combs := make([]rule.Combiner, len(c.CombReqs))
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
	Idx              string `json:"-"`
}

// ID func implement Ruler interface
func (s *ServerCfg) ID() string { return s.Idx }

// String func implement Ruler interface
func (s *ServerCfg) String() string {
	return fmt.Sprintf("ServerCfg id: %s, prefix: %s", s.Idx, s.PPrefix)
}

// SetID func implement Ruler interface
func (s *ServerCfg) SetID(id string) { s.Idx = id }

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
	Idx    string `json:"-"`
}

// ID func implement Ruler interface
func (s *ReverseServerCfg) ID() string { return s.Idx }

// String func implement Ruler interface
func (s *ReverseServerCfg) String() string {
	return fmt.Sprintf("ReverseServerCfg id: %s, Addr: %s", s.Idx, s.AAddr)
}

// SetID func implement Ruler interface
func (s *ReverseServerCfg) SetID(id string) { s.Idx = id }

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
	Idx     string `json:"-"`
}

// ID func implement Ruler interface
func (c *CombineReqCfg) ID() string { return c.Idx }

// String func implement Ruler interface
func (c *CombineReqCfg) String() string {
	return fmt.Sprintf("CombineReqCfg id: %s, prefix: %s", c.Idx, c.PPath)
}

// SetID func implement Ruler interface
func (c *CombineReqCfg) SetID(id string) { c.Idx = id }

// ServerName () string
func (c *CombineReqCfg) ServerName() string { return c.SrvName }

// Path ...
func (c *CombineReqCfg) Path() string { return c.PPath }

// Field ...
func (c *CombineReqCfg) Field() string { return c.FField }

// Method ...
func (c *CombineReqCfg) Method() string { return c.MMethod }

// NocacheCfg ...
type NocacheCfg struct {
	Regexp string `json:"regular"`
	Idx    string `json:"-"`
}

// String ...
func (i *NocacheCfg) String() string { return fmt.Sprintf("NocacheCfg: %s", i.Regexp) }

// ID ...
func (i *NocacheCfg) ID() string { return i.Idx }

// SetID ...
func (i *NocacheCfg) SetID(id string) { i.Idx = id }

// Regular ...
func (i *NocacheCfg) Regular() string { return i.Regexp }
