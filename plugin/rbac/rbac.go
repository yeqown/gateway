package rbac

import (
	"errors"
	"hash"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/plugin"
)

var (
	_               plugin.Plugin = &RBAC{}
	errNoPermission               = errors.New("No Permission")
)

// New ... only return a RBAC instance, and must load rules manually
func New(fieldName string, us []rule.User) *RBAC {
	if fieldName == "" {
		fieldName = "user_id"
	}
	r := &RBAC{
		enabled:     true,
		status:      plugin.Working,
		userIDField: fieldName,
	}
	r.LoadUsers(us)
	return r
}

// RBAC ...
type RBAC struct {
	enabled     bool
	status      plugin.PlgStatus
	userIDField string
	urlHashMap  map[string]rule.Permission
	md5er       hash.Hash

	users map[string]rule.User
	// roles       map[string]rule.Role
	// permissions map[string]rule.Permission
}

// LoadPermissions ...
// func (r *RBAC) LoadPermissions(perms []rule.Permission) {}
// LoadRoles ...
// func (r *RBAC) LoadRoles(roles []rule.Role) {}

// LoadUsers ...
func (r *RBAC) LoadUsers(users []rule.User) {
	r.users = make(map[string]rule.User)
	for _, u := range users {
		if _, ex := r.users[u.UserID()]; ex {
			panic("duplicated user_id")
		}
		r.users[u.UserID()] = u
	}
}

// LoadURLRules ...
func (r *RBAC) LoadURLRules(rules []rule.PermitURL) {
	r.urlHashMap = make(map[string]rule.Permission)
	for _, rule := range rules {
		hashed := r.hashURI(rule.URL())
		perm := rule.Permission()
		if perm == nil {
			panic("could not be nil Permission")
		}
		r.urlHashMap[hashed] = perm
	}
}

// Handle ....
func (r *RBAC) Handle(ctx *plugin.Context) {
	defer plugin.Recover("plugin.rbac")
	var (
		permitted bool
		need      bool
	)
	// permit url
	if permitted, need = r.permit(ctx.Path,
		ctx.Form.Get(r.userIDField)); need && !permitted {
		// 需要权限才能访问，且没有权限
		ctx.SetError(errNoPermission)
		ctx.Abort()
		return
	}

	ctx.Next()
}

// Name ...
func (r *RBAC) Name() string {
	return "plugin.rabc"
}

// Enabled ...
func (r *RBAC) Enabled() bool {
	return r.enabled
}

// Status ...
func (r *RBAC) Status() plugin.PlgStatus {
	return r.status
}

// Enable ...
func (r *RBAC) Enable(enabled bool) {
	r.enabled = enabled
	r.status = plugin.Working
	if !enabled {
		r.status = plugin.Stopped
	}
}

func (r *RBAC) permit(uri, userID string) (permitted, need bool) {
	hashed := r.hashURI(uri)
	perm, ex := r.urlHashMap[hashed]
	if !ex {
		// no need to permit the request
		need = false
		return
	}

	need = true
	if userID == "" {
		// userID is empty
		// TODO: support default role
		permitted = false
		return
	}
	user, ok := r.users[userID]
	if !ok {
		// missed userID
		permitted = false
		return
	}

	// brute force
	for _, role := range user.Roles() {
		if permitted := role.Permit(perm); permitted {
			break
		}
	}
	return
}

// func (r *RBAC) needToPermit(uri string) (perm, bool) {
// 	hashed := r.hashURI(uri)
// 	perm, ok := r.urlHashMap[hashed]
// 	return perm, ok
// }

func (r *RBAC) hashURI(uri string) string {
	r.md5er.Reset()
	_, err := r.md5er.Write([]byte(uri))
	if err != nil {
		panic(err)
	}

	return string(r.md5er.Sum(nil))
}
