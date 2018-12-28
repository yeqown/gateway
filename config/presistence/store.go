package presistence

import "github.com/yeqown/gateway/config/rule"

// PluginCode ...
type PluginCode byte

const (
	// PlgCodeProxyPath ...
	PlgCodeProxyPath PluginCode = 1 << iota
	// PlgCodeProxyServer ...
	PlgCodeProxyServer
	// PlgCodeProxyReverseSrv ...
	PlgCodeProxyReverseSrv
	// PlgCodeCache ...
	PlgCodeCache
	// PlgCodeRatelimit ...
	PlgCodeRatelimit
)

var (
	codes = []PluginCode{
		PlgCodeProxyPath,
		PlgCodeProxyServer,
		PlgCodeProxyReverseSrv,
		PlgCodeCache,
		PlgCodeRatelimit,
	}
)

// ListPlgByCode to get all changed plugin config
func ListPlgByCode(code PluginCode) []PluginCode {
	result := make([]PluginCode, 0, len(codes))
	for _, c := range codes {
		if (c & code) == c {
			result = append(result, c)
		}
	}
	return result
}

// ChangedChan ...
type ChangedChan struct {
	Code PluginCode
}

// Instance includes all config fields will be used
type Instance struct {
	ProxyServerRules    []rule.ServerRuler              `json:"server_rules"`
	ProxyPathRules      []rule.PathRuler                `json:"path_rules"`
	ProxyReverseServers map[string][]rule.ReverseServer `json:"reverse_servers"`
	Nocache             []rule.Nocacher                 `json:"nocache_rule"`
	Users               []rule.User                     `json:"users"`
}

// Store ... to add, del, query, update config rule ~
type Store interface {
	ServerRulerManager
	PathRulerManager
	ReverseServerManager
	NocacherManager

	PermissionManager
	RoleManager
	UserManager
	PermitURLManager

	// Instance get global config instance
	Instance() *Instance

	// Updated return a <-chan
	Updated() <-chan ChangedChan
}

// ServerRulerManager to contains all ServerRuler manage funcs
type ServerRulerManager interface {
	// NewServerRule func
	NewServerRule(r rule.ServerRuler) error
	// DelServerRule func
	DelServerRule(id string) error
	// UpdateServerRule func ...
	UpdateServerRule(id string, r rule.ServerRuler) error
	// PathRuleByID ...
	ServerRuleByID(id string) rule.ServerRuler
	// ServerRulesPage func
	ServerRulesPage(offset, limit int) ([]rule.ServerRuler, int)
}

// PathRulerManager rule.PathRuler manage funcs !!!
type PathRulerManager interface {
	// NewPathRule func
	NewPathRule(r rule.PathRuler) error
	// DelPathRule func
	DelPathRule(id string) error
	// UpdatePathRule func ...
	UpdatePathRule(id string, r rule.PathRuler) error
	// PathRuleByID ...
	PathRuleByID(id string) rule.PathRuler
	// PathRulesPage func
	PathRulesPage(offset, limit int) ([]rule.PathRuler, int)
}

// ReverseServerManager rule.ReverseServer manage funcs !!!
type ReverseServerManager interface {
	// NewReverseServer func
	NewReverseServer(group string, s rule.ReverseServer) error
	// DelReverseServer func
	DelReverseServer(id string) error
	// DelReverseServerGroup func
	DelReverseServerGroup(group string) error
	// UpdateReverseServerGroupName ...
	UpdateReverseServerGroupName(group string, newname string) error
	// UpdateReverseServer func ...
	UpdateReverseServer(id string, s rule.ReverseServer) error
	// ReverseServerByID ...
	ReverseServerByID(group string, id string) rule.ReverseServer
	// ReverseServerByGroup 根据分组来分页展示
	ReverseServerByGroup(group string, offset, limit int) ([]rule.ReverseServer, int)
	// ReverseServerGroups 获取分组及分组下配置计数
	ReverseServerGroups() map[string]int
}

// NocacherManager rule.Nocacher manage funcs !!!
type NocacherManager interface {
	// NewNocacheRule func
	NewNocacheRule(c rule.Nocacher) error
	// DelNocacheRule func
	DelNocacheRule(id string) error
	// UpdateNocacheRule func ...
	UpdateNocacheRule(id string, c rule.Nocacher) error
	// NocacheRules ...
	NocacheRules(offset, limit int) ([]rule.Nocacher, int)
	// NocacheRuleByID ...
	NocacheRuleByID(id string) rule.Nocacher
}

// PermissionManager rule.Permission manage funcs
type PermissionManager interface {
	NewPermission(r rule.Permission) error
	DelPermission(id string) error
	EditPermission(id string, r rule.Permission) error
	PermissionPage(limit, offset int) ([]rule.Permission, int)
}

// RoleManager rule.Role manage funcs
type RoleManager interface {
	NewRole(r rule.Role) error
	DelRole(id string) error
	EditRole(id string, r rule.Role) error
	RolePage(limit, offset int) ([]rule.Role, int)
	AssignPerm(id string, permids ...string) error
	RevokePerm(id string, permids ...string) error
}

// UserManager rule.User manage funcs
type UserManager interface {
	NewUser(r rule.User) error
	DelUser(id string) error
	EditUser(id string, r rule.User) error
	UserPage(limit, offset int) ([]rule.User, int)
	AssignRole(id, roleid string) error
	RevokeRole(id, roleid string) error
}

// PermitURLManager rule.PermitURL manage funcs
type PermitURLManager interface {
	NewPermitURL(r rule.PermitURL) error
	DelPermitURL(id string) error
	EditPermitURL(id string, r rule.PermitURL) error
	PermitURLPage(limit, offset int) ([]rule.PermitURL, int)
}
