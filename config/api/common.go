package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
	"gopkg.in/go-playground/validator.v9"
)

var (
	global   presistence.Store
	decoder  *schema.Decoder
	validate *validator.Validate

	_ rule.PathRuler   = &apiPathRuler{}
	_ rule.ServerRuler = &apiServerRuler{}
	_ rule.Nocacher    = &apiNocacher{}
	_ rule.Combiner    = &apiCombReqCfg{}

	_ rule.Permission = &apiPermission{}
	_ rule.Role       = &apiRole{}
	_ rule.PermitURL  = &apiPermitURL{}
	_ rule.User       = &apiUser{}
)

func init() {
	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)
	decoder.SetAliasTag("form")

	validate = validator.New()
	validate.SetTagName("valid")
}

// Global ...
func Global() presistence.Store {
	if global == nil {
		panic("server init wrong, global varible is nil")
	}
	return global
}

// SetGlobal ..
func SetGlobal(store presistence.Store) {
	global = store
}

// Bind form
func Bind(dst interface{}, req *http.Request) error {
	form := utils.ParseRequestForm(req)
	return decoder.Decode(dst, form)
}

// Valid ...
func Valid(v interface{}) error {
	err := validate.Struct(v)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errmsg := ""
		for _, fieldErr := range errs {
			errmsg += fmt.Sprintf("表单%s不符合%s校验要求;",
				fieldErr.Field(), fieldErr.Tag(),
			)
		}
		return fmt.Errorf("参数非法: %s", errmsg)

	}
	return err
}

// ResponseWithError ...
func ResponseWithError(w http.ResponseWriter, resp interface{}, err error) {
	code.FillCodeInfo(resp,
		code.NewCodeInfo(code.CodeParamInvalid, err.Error()))
	utils.ResponseJSON(w, resp)
	return
}

type commonResp struct {
	code.CodeInfo
}

// !!!! 以下的所有结构体定义和函数都是用于声明package.rule中的interface !!!!
// 如：Nocaher, ServerRuler, PathRuler等等
// 用于API表单使用, 同时定义Interface到response struct的转换

type apiPathRuler struct {
	PPath        string           `json:"path" form:"path" valid:"required"`
	RRewritePath string           `json:"rewrite_path" form:"rewrite_path" valid:"required"`
	MMethod      string           `json:"method" form:"method" valid:"required"`
	SrvName      string           `json:"server_name" form:"server_name" valid:"required"`
	CombReqs     []*apiCombReqCfg `json:"combine_req_cfgs" form:"combine_req_cfgs"`
	NeedComb     bool             `json:"need_combine" form:"need_combine"`
	Idx          string           `json:"id" form:"-"`
}

func (c *apiPathRuler) ID() string { return c.Idx }
func (c *apiPathRuler) String() string {
	return fmt.Sprintf("%v", *c)
}
func (c *apiPathRuler) SetID(id string)     { c.Idx = id }
func (c *apiPathRuler) Path() string        { return c.PPath }
func (c *apiPathRuler) Method() string      { return c.MMethod }
func (c *apiPathRuler) ServerName() string  { return c.SrvName }
func (c *apiPathRuler) RewritePath() string { return c.RRewritePath }
func (c *apiPathRuler) NeedCombine() bool   { return c.NeedComb }
func (c *apiPathRuler) CombineReqCfgs() []rule.Combiner {
	combs := make([]rule.Combiner, len(c.CombReqs))
	for idx, c := range c.CombReqs {
		combs[idx] = c
	}
	return combs
}
func loadFromPathRuler(r rule.PathRuler) *apiPathRuler {
	reqs := r.CombineReqCfgs()
	combReqs := make([]*apiCombReqCfg, len(reqs))

	for idx, r := range reqs {
		combReqs[idx] = loadFromCombiner(r)
	}

	return &apiPathRuler{
		PPath:        r.Path(),
		RRewritePath: r.RewritePath(),
		MMethod:      r.Method(),
		SrvName:      r.ServerName(),
		CombReqs:     combReqs,
		NeedComb:     r.NeedCombine(),
		Idx:          r.ID(),
	}
}

type apiServerRuler struct {
	PPrefix          string `json:"prefix" form:"prefix" valid:"required"`
	SServerName      string `json:"server_name" form:"server_name" valid:"required"`
	NNeedStripPrefix bool   `json:"need_strip_prefix" form:"need_strip_prefix"`
	Idx              string `json:"id" form:"-"`
}

func (s *apiServerRuler) ID() string      { return s.Idx }
func (s *apiServerRuler) SetID(id string) { s.Idx = id }
func (s *apiServerRuler) String() string {
	return fmt.Sprintf("ServerCfg id: %s, prefix: %s", s.Idx, s.PPrefix)
}
func (s *apiServerRuler) Prefix() string        { return s.PPrefix }
func (s *apiServerRuler) ServerName() string    { return s.SServerName }
func (s *apiServerRuler) NeedStripPrefix() bool { return s.NNeedStripPrefix }
func loadFromServerRuler(r rule.ServerRuler) *apiServerRuler {
	return &apiServerRuler{
		PPrefix:          r.Prefix(),
		SServerName:      r.ServerName(),
		NNeedStripPrefix: r.NeedStripPrefix(),
		Idx:              r.ID(),
	}
}

type apiReverseSrver struct {
	NName  string `json:"name" form:"name" valid:"required"`
	AAddr  string `json:"addr" form:"addr" valid:"required"`
	Weight int    `json:"weight" form:"weight" valid:"required"`
	GGroup string `json:"group" form:"group" valid:"required"`
	Idx    string `json:"id" form:"-"`
}

func (s *apiReverseSrver) ID() string      { return s.Idx }
func (s *apiReverseSrver) SetID(id string) { s.Idx = id }
func (s *apiReverseSrver) String() string {
	return fmt.Sprintf("apiReverseSrver id: %s, Addr: %s", s.Idx, s.AAddr)
}
func (s *apiReverseSrver) Group() string { return s.GGroup }
func (s *apiReverseSrver) Name() string  { return s.NName }
func (s *apiReverseSrver) Addr() string  { return s.AAddr }
func (s *apiReverseSrver) W() int        { return s.Weight }
func loadFromReverseServer(r rule.ReverseServer) *apiReverseSrver {
	return &apiReverseSrver{
		NName:  r.Name(),
		AAddr:  r.Addr(),
		Weight: r.W(),
		GGroup: r.Group(),
		Idx:    r.ID(),
	}
}

type apiCombReqCfg struct {
	SrvName string `json:"server_name" form:"server_name" valid:"required"` // http://ip:port/path?params
	PPath   string `json:"path" form:"path" valid:"required"`               // path `/request/path`
	FField  string `json:"field" form:"field" valid:"required"`             // want got field
	MMethod string `json:"method" form:"method" valid:"required"`           // want match method
	Idx     string `json:"id" form:"-"`
}

func (c *apiCombReqCfg) ID() string      { return c.Idx }
func (c *apiCombReqCfg) SetID(id string) { c.Idx = id }
func (c *apiCombReqCfg) String() string {
	return fmt.Sprintf("apiCombReqCfg id: %s, prefix: %s", c.Idx, c.PPath)
}
func (c *apiCombReqCfg) ServerName() string { return c.SrvName }
func (c *apiCombReqCfg) Path() string       { return c.PPath }
func (c *apiCombReqCfg) Field() string      { return c.FField }
func (c *apiCombReqCfg) Method() string     { return c.MMethod }
func loadFromCombiner(r rule.Combiner) *apiCombReqCfg {
	return &apiCombReqCfg{
		SrvName: r.ServerName(),
		PPath:   r.Path(),
		FField:  r.Field(),
		MMethod: r.Method(),
		Idx:     r.ID(),
	}
}

type apiNocacher struct {
	Regexp   string `json:"regular" form:"regular" valid:"required"`
	EEnabled bool   `json:"enabled" form:"enabled"`
	Idx      string `json:"id" form:"-"`
}

func (i *apiNocacher) String() string  { return fmt.Sprintf("apiNocacher: %s", i.Regexp) }
func (i *apiNocacher) ID() string      { return i.Idx }
func (i *apiNocacher) SetID(id string) { i.Idx = id }
func (i *apiNocacher) Regular() string { return i.Regexp }
func (i *apiNocacher) Enabled() bool   { return i.EEnabled }
func loadFromNocacher(r rule.Nocacher) *apiNocacher {
	return &apiNocacher{
		Regexp:   r.Regular(),
		Idx:      r.ID(),
		EEnabled: r.Enabled(),
	}
}

// rbac support structs

type apiPermission struct {
	Idx       string `json:"id" form:"-" valid:"-"`
	AAction   string `json:"action" form:"action" valid:"required"`
	RResource string `json:"resource" form:"resource" valid:"required"`
}

func (p *apiPermission) ID() string                       { return p.Idx }
func (p *apiPermission) SetID(id string)                  { p.Idx = id }
func (p *apiPermission) Action() string                   { return p.AAction }
func (p *apiPermission) Resource() string                 { return p.RResource }
func (p *apiPermission) Match(other rule.Permission) bool { panic("not Implemented") }
func loadFromPermission(r rule.Permission) *apiPermission {
	return &apiPermission{
		Idx:       r.ID(),
		AAction:   r.Action(),
		RResource: r.Resource(),
	}
}

type apiRole struct {
	Idx     string           `json:"id" form:"-" valid:"-"`
	Perms   []*apiPermission `json:"permissions" form:"-" valid:"-"`
	PermIDs []string         `json:"-" form:"perm_ids" valid:"required"`
	NName   string           `json:"name" form:"name" valid:"required"`
}

func (p *apiRole) ID() string                            { return p.Idx }
func (p *apiRole) SetID(id string)                       { p.Idx = id }
func (p *apiRole) Name() string                          { return p.NName }
func (p *apiRole) Permit(perm rule.Permission) bool      { panic("not Implemented") }
func (p *apiRole) Assign(perms ...rule.Permission) error { panic("not Implemented") }
func (p *apiRole) Revoke(prem rule.Permission) error     { panic("not Implemented") }
func (p *apiRole) Permissions() []rule.Permission {
	// log.Printf("permissions: %v", *p)
	rules := make([]rule.Permission, len(p.PermIDs))
	for idx, permID := range p.PermIDs {
		rules[idx] = &apiPermission{
			Idx: permID,
		}
	}
	return rules
}
func loadFormRole(r rule.Role) *apiRole {
	// log.Printf("api load role: %v", r)
	permissions := r.Permissions()
	perms := make([]*apiPermission, len(permissions))
	for idx, perm := range permissions {
		perms[idx] = loadFromPermission(perm)
	}
	return &apiRole{
		Idx:   r.ID(),
		Perms: perms,
		NName: r.Name(),
	}
}

type apiUser struct {
	Idx     string     `json:"id" form:"-" valid:"-"`
	UUserID string     `json:"user_id" form:"user_id" valid:"required"`
	RoleIDs []string   `json:"-" form:"role_ids" valid:"required"`
	RRoles  []*apiRole `json:"roles" form:"-" valid:"-"`
}

func (p *apiUser) ID() string                      { return p.Idx }
func (p *apiUser) SetID(id string)                 { p.Idx = id }
func (p *apiUser) UserID() string                  { return p.UUserID }
func (p *apiUser) Assign(roles ...rule.Role) error { panic("not Implemented") }
func (p *apiUser) Revoke(role rule.Role) error     { panic("not Implemented") }
func (p *apiUser) Roles() []rule.Role {
	rules := make([]rule.Role, len(p.RoleIDs))
	for idx, roleID := range p.RoleIDs {
		rules[idx] = &apiRole{
			Idx: roleID,
		}
	}
	return rules
}
func loadFormUser(r rule.User) *apiUser {
	oriRoles := r.Roles()
	roles := make([]*apiRole, len(oriRoles))
	for idx, role := range oriRoles {
		roles[idx] = loadFormRole(role)
	}
	return &apiUser{
		Idx:     r.ID(),
		UUserID: r.UserID(),
		RRoles:  roles,
	}
}

type apiPermitURL struct {
	Idx    string         `json:"id" form:"id" valid:"-"`
	URI    string         `json:"uri" form:"uri" valid:"required"`
	Perm   *apiPermission `json:"permission" form:"-"`
	PermID string         `json:"-" form:"perm_id" valid:"required"`
}

func (p *apiPermitURL) ID() string      { return p.Idx }
func (p *apiPermitURL) SetID(id string) { p.Idx = id }
func (p *apiPermitURL) URL() string     { return p.URI }
func (p *apiPermitURL) Permission() rule.Permission {
	return &apiPermission{
		Idx: p.PermID,
	}
}
func loadFormPermitURL(r rule.PermitURL) *apiPermitURL {
	return &apiPermitURL{
		Idx:    r.ID(),
		URI:    r.URL(),
		PermID: r.Permission().ID(),
		Perm:   loadFromPermission(r.Permission()),
	}
}
