package rule

// Permission 定义权限的具体表述
type Permission interface {
	Ruler
	Action() string
	Resource() string
	Match(p Permission) bool
}

// Role 接口定义
type Role interface {
	Ruler
	Name() string
	Permissions() []Permission
	Permit(perm Permission) bool
	Assign(perms ...Permission) error
	Revoke(prem Permission) error
}

// User 接口定义用户数据结构的基本方法 map[user_id]
type User interface {
	Ruler
	UserID() string
	Roles() []Role
	Assign(roles ...Role) error
	Revoke(role Role) error
}

// PermitURL 需要鉴权的请求路径
type PermitURL interface {
	Ruler
	URL() string
	Permission() Permission
}
