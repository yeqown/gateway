package rule

import "fmt"

// Nocacher ...
type Nocacher interface {
	Ruler
	Regular() string
	// SetID(string)
}

// NocacheManager ...
type NocacheManager interface {
	GetByID(id string) Nocacher
	GetSlice(start, end int) []Nocacher
	Count() int
}

// NocacheCfg ...
type NocacheCfg struct {
	Regexp string `json:"regular"`
	id     string
}

// Regular ...
func (i *NocacheCfg) Regular() string {
	return i.Regexp
}

// String ...
func (i *NocacheCfg) String() string {
	return fmt.Sprintf("Nocacher: %s", i.Regexp)
}

// ID ...
func (i *NocacheCfg) ID() string {
	return i.id
}

// NocacheManagerInstance ...
type NocacheManagerInstance struct {
	rules []Nocacher
}

// Count ... func to count all
func (m *NocacheManagerInstance) Count() int {
	return len(m.rules)
}

// GetByID ...
func (m *NocacheManagerInstance) GetByID(id string) Nocacher {
	for _, rule := range m.rules {
		if rule.ID() == id {
			return rule
		}
	}
	return nil
}

// GetSlice ...
func (m *NocacheManagerInstance) GetSlice(start, end int) []Nocacher {
	return m.rules[start:end]
}
