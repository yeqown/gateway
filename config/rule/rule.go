package rule

// Ruler interface
type Ruler interface {
	ID() string
	String() string
}

type defaultRuler struct{}

// ID ...
func ID() string { return "NotImplement" }

// String ...
func String() string { return "NotImplement" }
