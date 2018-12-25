package rule

// Ruler interface
type Ruler interface {
	ID() string
	SetID(id string)
	// String() string
}

// Nocacher ...
type Nocacher interface {
	Ruler
	Regular() string
	Enabled() bool
}

// PathRuler 用于单个url配置
type PathRuler interface {
	Ruler
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
	Ruler
	Prefix() string
	ServerName() string // return group name
	NeedStripPrefix() bool
}

// ReverseServer 单向代理服务器配置
type ReverseServer interface {
	Ruler
	Name() string
	Addr() string
	Group() string
	W() int
}

// Combiner 用于合并请求时候用的配置
type Combiner interface {
	Ruler
	ServerName() string
	Path() string
	Field() string
	Method() string
}
