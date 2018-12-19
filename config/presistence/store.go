package presistence

import "github.com/yeqown/gateway/config/rule"

// Store ... to add, del, query, update config rule ~
type Store interface {
	// Add op collections
	// NewServerRule func
	NewServerRule(r rule.ServerRuler) error
	// NewPathRule func
	NewPathRule(r rule.PathRuler) error

	//TODO: NewReverseServerGroup()

	// NewReverseServer func
	NewReverseServer(group string, s rule.ReverseServer) error
	// NewNocacheRule func
	NewNocacheRule(c rule.Nocacher) error

	// Query op collecitons
	// Instance get global config instance
	Instance() *Instance

	// NocacheRules ...
	NocacheRules(start, end int) []rule.Nocacher
	// NocacheRuleByID ...
	NocacheRuleByID(id string) rule.Nocacher
	// NocacheRulesCount ...
	NocacheRulesCount() int

	// PathRuleByID ...
	ServerRuleByID(id string) rule.ServerRuler
	// ServerRulesPage func
	ServerRulesPage(start, end int) []rule.ServerRuler
	// ServerRulesCount func
	ServerRulesCount() int

	// PathRuleByID ...
	PathRuleByID(id string) rule.PathRuler
	// PathRulesPage func
	PathRulesPage(start, end int) []rule.PathRuler
	// PathRulesCount func
	PathRulesCount() int

	// ReverseServerByID ...
	ReverseServerByID(group string, id string) rule.ReverseServer
	// ReverseServerGroup func
	ReverseServerGroup(group string, start, end int) []rule.ReverseServer
	// ReverseServerGroupCount func
	ReverseServerGroupCount() int
	// ReverseServerGroupPageCount func
	ReverseServerGroupPageCount(group string) int

	// Update op collections
	// UpdateGateConfig func ...
	UpdateGateConfig(logpath string, port int) error
	// UpdateServerRule func ...
	UpdateServerRule(id string, r rule.ServerRuler) error
	// UpdatePathRule func ...
	UpdatePathRule(id string, r rule.PathRuler) error
	// UpdateReverseServer func ...
	UpdateReverseServer(id string, s rule.ReverseServer) error
	// UpdateNocacheRule func ...
	UpdateNocacheRule(id string, c rule.Nocacher) error

	// Del op collections
	// DelServerRule func
	DelServerRule(id string) error
	// DelPathRule func
	DelPathRule(id string) error
	// DelReverseServer func
	DelReverseServer(id string) error
	// DelReverseServerGroup func
	DelReverseServerGroup(group string) error
	// DelNocacheRule func
	DelNocacheRule(id string) error
}

// Instance includes all config fields will be used
type Instance struct {
	Logpath             string                          `json:"logpath"`
	Port                int                             `json:"port"`
	ProxyServerRules    []rule.ServerRuler              `json:"server_rules"`
	ProxyPathRules      []rule.PathRuler                `json:"path_rules"`
	ProxyReverseServers map[string][]rule.ReverseServer `json:"reverse_servers"`
	Nocache             []rule.Nocacher                 `json:"nocache_rule"`
}
