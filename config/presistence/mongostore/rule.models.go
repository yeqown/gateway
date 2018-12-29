package mongostore

import (
	"github.com/yeqown/gateway/config/rule"
	"gopkg.in/mgo.v2/bson"
)

var (
	_ rule.Nocacher      = &nocacherModel{}
	_ rule.PathRuler     = &pathRulerModel{}
	_ rule.Combiner      = &combinerModel{}
	_ rule.ServerRuler   = &serverRulerModel{}
	_ rule.ReverseServer = &reverseServerModel{}

	_ rule.Permission = &permissionModel{}
	_ rule.Role       = &roleModel{}
	_ rule.User       = &userModel{}
	_ rule.PermitURL  = &permitURLModel{}
)

// structs for plugin.proxy
type nocacherModel struct {
	Idx      bson.ObjectId `bson:"_id"`
	Regexp   string        `bson:"regular"`
	EEnabled bool          `bson:"enabled"`
}

func (c *nocacherModel) Regular() string { return c.Regexp }
func (c *nocacherModel) Enabled() bool   { return c.EEnabled }
func (c *nocacherModel) ID() string      { return c.Idx.Hex() }
func (c *nocacherModel) SetID(id string) { c.Idx = bson.ObjectIdHex(id) }
func loadNocacherModelFromNocacher(r rule.Nocacher) *nocacherModel {
	return &nocacherModel{
		Idx:      bson.ObjectIdHex(r.ID()),
		Regexp:   r.Regular(),
		EEnabled: r.Enabled(),
	}
}

type combinerModel struct {
	Idx         bson.ObjectId `bson:"_id"`
	SServerName string        `bson:"sever_name"`
	PPath       string        `bson:"path"`
	FField      string        `bson:"field"`
	MMethod     string        `bson:"method"`
}

func (c *combinerModel) ID() string         { return c.Idx.Hex() }
func (c *combinerModel) SetID(id string)    { c.Idx = bson.ObjectIdHex(id) }
func (c *combinerModel) ServerName() string { return c.SServerName }
func (c *combinerModel) Path() string       { return c.PPath }
func (c *combinerModel) Field() string      { return c.FField }
func (c *combinerModel) Method() string     { return c.MMethod }
func loadCombinerModelFromCombiner(r rule.Combiner) *combinerModel {
	if len(r.ID()) == 0 {
		r.SetID(bson.NewObjectId().Hex())
	}
	return &combinerModel{
		Idx:         bson.ObjectIdHex(r.ID()),
		SServerName: r.ServerName(),
		PPath:       r.Path(),
		FField:      r.Field(),
		MMethod:     r.Method(),
	}
}

type pathRulerModel struct {
	Idx             bson.ObjectId    `bson:"_id"`
	PPath           string           `bson:"path"`
	MMethod         string           `bson:"method"`
	SServerName     string           `bson:"server_name"`
	RRewritePath    string           `bson:"rewrite_path"`
	NNeedCombine    bool             `bson:"need_combine"`
	CCombineReqCfgs []*combinerModel `bson:"combine_req_cfgs"`
}

func (c *pathRulerModel) ID() string          { return c.Idx.Hex() }
func (c *pathRulerModel) SetID(id string)     { c.Idx = bson.ObjectIdHex(id) }
func (c *pathRulerModel) Path() string        { return c.PPath }
func (c *pathRulerModel) Method() string      { return c.MMethod }
func (c *pathRulerModel) ServerName() string  { return c.SServerName }
func (c *pathRulerModel) RewritePath() string { return c.RRewritePath }
func (c *pathRulerModel) NeedCombine() bool   { return c.NNeedCombine }
func (c *pathRulerModel) CombineReqCfgs() []rule.Combiner {
	rules := make([]rule.Combiner, len(c.CCombineReqCfgs))
	for idx, r := range c.CCombineReqCfgs {
		rules[idx] = r
	}
	return rules
}
func loadPathRulerModelFromPathRuler(r rule.PathRuler) *pathRulerModel {
	combReqs := make([]*combinerModel, len(r.CombineReqCfgs()))
	for idx, r := range r.CombineReqCfgs() {
		combReqs[idx] = loadCombinerModelFromCombiner(r)
	}

	return &pathRulerModel{
		Idx:             bson.ObjectIdHex(r.ID()),
		PPath:           r.Path(),
		MMethod:         r.Method(),
		SServerName:     r.ServerName(),
		RRewritePath:    r.RewritePath(),
		NNeedCombine:    r.NeedCombine(),
		CCombineReqCfgs: combReqs,
	}
}

type serverRulerModel struct {
	Idx              bson.ObjectId `bson:"_id"`
	PPrefix          string        `bson:"prefix"`
	SServerName      string        `bson:"server_name"`
	NNeedStripPrefix bool          `bson:"need_strip_prefix"`
}

func (c *serverRulerModel) ID() string            { return c.Idx.Hex() }
func (c *serverRulerModel) SetID(id string)       { c.Idx = bson.ObjectIdHex(id) }
func (c *serverRulerModel) Prefix() string        { return c.PPrefix }
func (c *serverRulerModel) ServerName() string    { return c.SServerName }
func (c *serverRulerModel) NeedStripPrefix() bool { return c.NNeedStripPrefix }
func loadServerRulerModelFromServerRuler(r rule.ServerRuler) *serverRulerModel {
	return &serverRulerModel{
		Idx:              bson.ObjectIdHex(r.ID()),
		PPrefix:          r.Prefix(),
		SServerName:      r.ServerName(),
		NNeedStripPrefix: r.NeedStripPrefix(),
	}
}

type reverseServerModel struct {
	Idx    bson.ObjectId `bson:"_id"`
	NName  string        `bson:"name"`
	AAddr  string        `bson:"addr"`
	WW     int           `bson:"w"`
	GGroup string        `bson:"group"`
}

func (c *reverseServerModel) ID() string      { return c.Idx.Hex() }
func (c *reverseServerModel) SetID(id string) { c.Idx = bson.ObjectIdHex(id) }
func (c *reverseServerModel) Name() string    { return c.NName }
func (c *reverseServerModel) Addr() string    { return c.AAddr }
func (c *reverseServerModel) W() int          { return c.WW }
func (c *reverseServerModel) Group() string   { return c.GGroup }
func loadReverseServerModelFromReverseServer(r rule.ReverseServer) *reverseServerModel {
	return &reverseServerModel{
		Idx:    bson.ObjectIdHex(r.ID()),
		NName:  r.Name(),
		AAddr:  r.Addr(),
		WW:     r.W(),
		GGroup: r.Group(),
	}
}

// struct for RBAC plugins

type permissionModel struct {
	Idx       bson.ObjectId `bson:"_id"`
	AAction   string        `bson:"action"`
	RResource string        `bson:"resource"`
}

func (p *permissionModel) ID() string       { return p.Idx.Hex() }
func (p *permissionModel) SetID(id string)  { p.Idx = bson.ObjectIdHex(id) }
func (p *permissionModel) Action() string   { return p.AAction }
func (p *permissionModel) Resource() string { return p.RResource }
func (p *permissionModel) Match(other rule.Permission) bool {
	if p.ID() == other.ID() {
		return true
	}
	return p.AAction == other.Action() && p.RResource == other.Resource()
}
func loadPermissionModelFromPermission(r rule.Permission) *permissionModel {
	return &permissionModel{
		Idx:       bson.ObjectIdHex(r.ID()),
		AAction:   r.Action(),
		RResource: r.Resource(),
	}
}

type roleModel struct {
	Idx         bson.ObjectId              `bson:"_id"`
	PermIDs     []bson.ObjectId            `bson:"perm_ids"`
	permissions map[string]rule.Permission `bson:"-"`
	NName       string                     `bson:"name"`
}

func (r *roleModel) ID() string      { return r.Idx.Hex() }
func (r *roleModel) SetID(id string) { r.Idx = bson.ObjectIdHex(id) }
func (r *roleModel) Name() string    { return r.NName }
func (r *roleModel) Permissions() []rule.Permission {
	perms := make([]rule.Permission, len(r.permissions))
	counter := 0
	for _, v := range r.permissions {
		perms[counter] = v
		counter++
	}
	return perms
}
func (r *roleModel) Assign(perms ...rule.Permission) error {
	for _, perm := range perms {
		if _, ex := r.permissions[perm.ID()]; ex {
			continue
		}
		r.permissions[perm.ID()] = perm
		r.PermIDs = append(r.PermIDs, bson.ObjectIdHex(perm.ID()))
	}
	return nil
}
func (r *roleModel) Revoke(perm rule.Permission) error {
	delete(r.permissions, perm.ID())

	delIdx := -1
	for idx, v := range r.PermIDs {
		if v.Hex() == perm.ID() {
			delIdx = idx
			break
		}
	}
	if delIdx == -1 || delIdx > len(r.PermIDs) {
		return nil
	}
	r.PermIDs = append(r.PermIDs[:delIdx], r.PermIDs[delIdx+1:]...)
	return nil
}
func (r *roleModel) Permit(perm rule.Permission) bool {
	var rslt bool
	for _, p := range r.permissions {
		if rslt = p.Match(perm); rslt {
			break
		}
	}
	// log.Infof("result: %v", rslt)
	return rslt
}
func loadRoleModelFromRole(r rule.Role) *roleModel {
	// log.Printf("%v, %v", r, r.Permissions())
	oriPerms := r.Permissions()
	ids := make([]bson.ObjectId, len(oriPerms))
	perms := make(map[string]rule.Permission)
	for idx, perm := range r.Permissions() {
		ids[idx] = bson.ObjectIdHex(perm.ID())
		perms[perm.ID()] = perm
	}

	return &roleModel{
		Idx:         bson.ObjectIdHex(r.ID()),
		permissions: perms,
		PermIDs:     ids,
		NName:       r.Name(),
	}
}

type userModel struct {
	Idx      bson.ObjectId        `bson:"_id"`
	UUserID  string               `bson:"user_id"`
	RoleIDs  []bson.ObjectId      `bson:"role_ids"`
	RolesMap map[string]rule.Role `bson:"-"`
}

func (u *userModel) ID() string      { return u.Idx.Hex() }
func (u *userModel) SetID(id string) { u.Idx = bson.ObjectIdHex(id) }
func (u *userModel) UserID() string  { return u.UUserID }
func (u *userModel) Roles() []rule.Role {
	roles := make([]rule.Role, len(u.RolesMap))
	counter := 0
	for _, v := range u.RolesMap {
		roles[counter] = v
		counter++
	}
	return roles
}

func (u *userModel) Assign(roles ...rule.Role) error {
	for _, role := range roles {
		if _, ex := u.RolesMap[role.ID()]; ex {
			continue
		}
		u.RolesMap[role.ID()] = role
		u.RoleIDs = append(u.RoleIDs, bson.ObjectIdHex(role.ID()))
	}
	return nil
}

func (u *userModel) Revoke(role rule.Role) error {
	delete(u.RolesMap, role.ID())

	delIdx := -1
	for idx, roleID := range u.RoleIDs {
		if roleID.Hex() == role.ID() {
			delIdx = idx
			break
		}
	}
	if delIdx == -1 || delIdx > len(u.RoleIDs) {
		return nil
	}
	u.RoleIDs = append(u.RoleIDs[:delIdx], u.RoleIDs[delIdx+1:]...)
	return nil
}

// func (u *userModel) Permit(perm rule.Permission) bool {}
func loadUserModelFromUser(r rule.User) *userModel {
	roles := r.Roles()
	roleIDs := make([]bson.ObjectId, len(roles))
	roleMap := make(map[string]rule.Role)

	for idx, role := range roles {
		roleIDs[idx] = bson.ObjectIdHex(role.ID())
		roleMap[role.ID()] = loadRoleModelFromRole(role)
	}

	return &userModel{
		Idx:      bson.ObjectIdHex(r.ID()),
		UUserID:  r.UserID(),
		RoleIDs:  roleIDs,
		RolesMap: roleMap,
	}
}

type permitURLModel struct {
	Idx    bson.ObjectId   `bson:"_id"`
	PermID bson.ObjectId   `bson:"perm_id"`
	Perm   rule.Permission `bson:"-"`
	URI    string          `bson:"url"`
}

func (u *permitURLModel) ID() string                  { return u.Idx.Hex() }
func (u *permitURLModel) SetID(id string)             { u.Idx = bson.ObjectIdHex(id) }
func (u *permitURLModel) URL() string                 { return u.URI }
func (u *permitURLModel) Permission() rule.Permission { return u.Perm }
func loadPermitURLModelFormPermitURL(r rule.PermitURL) *permitURLModel {
	permID := bson.ObjectIdHex(r.Permission().ID())
	return &permitURLModel{
		Idx:    bson.ObjectIdHex(r.ID()),
		PermID: permID,
		Perm:   r.Permission(),
		URI:    r.URL(),
	}
}
