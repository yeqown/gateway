package presistence

import "github.com/yeqown/gateway/config/rule"

// Store ... to add, del, query, update config rule ~
type Store interface {
	// TODO: Add
	// NewServerRule()
	// NewPathRule()
	// NewReverseServer()
	// NewNocacheRule()

	// Query ...
	Instance() *Instance

	// TODO: Update
	// UpdateGateConfig()
	// UpdateServerRule()
	// UpdatePathRule()
	// UpdateReverseServer()
	// UpdateNocacheRule()

	// TODO: Del
	// DelServerRule()
	// DelPathRule()
	// DelReverseServer()
	// DelNocacheRule()
}

// Instance includes all config fields will be used
type Instance struct {
	Logpath             string
	Port                int
	ProxyServerRules    []rule.ServerRuler
	ProxyPathRules      []rule.PathRuler
	ProxyReverseServers map[string][]rule.ReverseServer
	Nocache             []rule.Nocacher
}
