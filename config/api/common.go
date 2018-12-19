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

// bind form
func bind(dst interface{}, req *http.Request) error {
	form := utils.ParseRequestForm(req)
	return decoder.Decode(dst, form)
}

// valid ...
func valid(v interface{}) error {
	err := validate.Struct(v)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errmsg := ""
		for _, fieldErr := range errs {
			errmsg += fmt.Sprintf("Field_%s is invalid, validated by %s;",
				fieldErr.Field(), fieldErr.Tag(),
			)
		}
		return fmt.Errorf("validate error: %s", errmsg)

	}
	return err
}

// responseWithError ...
func responseWithError(w http.ResponseWriter, resp interface{}, err error) {
	code.FillCodeInfo(resp,
		code.NewCodeInfo(code.CodeParamInvalid, err.Error()))
	utils.ResponseJSON(w, resp)
	return
}

type commonResp struct {
	code.CodeInfo
}

// !!!! 以下的所有结构体定义和函数都是用于声明rule包中的interface !!!!
// 如：Nocaher, ServerRuler, PathRulerd等等
// 使用于API 表单定义方便使用

type formPathRuler struct {
	PPath        string            `json:"path" form:"path" valid:"required"`
	RRewritePath string            `json:"rewrite_path" form:"rewrite_path" valid:"required"`
	MMethod      string            `json:"method" form:"method" valid:"required"`
	SrvName      string            `json:"server_name" form:"server_name" valid:"required"`
	CombReqs     []*formCombReqCfg `json:"combine_req_cfgs" form:"combine_req_cfgs"`
	NeedComb     bool              `json:"need_combine" form:"need_combine"`
	Idx          string            `json:"id" form:"-"`
}

func (c *formPathRuler) ID() string { return c.Idx }
func (c *formPathRuler) String() string {
	return fmt.Sprintf("%v", *c)
}
func (c *formPathRuler) SetID(id string)     { c.Idx = id }
func (c *formPathRuler) Path() string        { return c.PPath }
func (c *formPathRuler) Method() string      { return c.MMethod }
func (c *formPathRuler) ServerName() string  { return c.SrvName }
func (c *formPathRuler) RewritePath() string { return c.RRewritePath }
func (c *formPathRuler) NeedCombine() bool   { return c.NeedComb }
func (c *formPathRuler) CombineReqCfgs() []rule.Combiner {
	combs := make([]rule.Combiner, len(c.CombReqs))
	for idx, c := range c.CombReqs {
		combs[idx] = c
	}
	return combs
}

type formServerRuler struct {
	PPrefix          string `json:"prefix" form:"prefix" valid:"required"`
	SServerName      string `json:"server_name" form:"server_name" valid:"required"`
	NNeedStripPrefix bool   `json:"need_strip_prefix" form:"need_strip_prefix"`
	Idx              string `json:"id" form:"-"`
}

func (s *formServerRuler) ID() string      { return s.Idx }
func (s *formServerRuler) SetID(id string) { s.Idx = id }
func (s *formServerRuler) String() string {
	return fmt.Sprintf("ServerCfg id: %s, prefix: %s", s.Idx, s.PPrefix)
}
func (s *formServerRuler) Prefix() string        { return s.PPrefix }
func (s *formServerRuler) ServerName() string    { return s.SServerName }
func (s *formServerRuler) NeedStripPrefix() bool { return s.NNeedStripPrefix }

type formReverseSrver struct {
	NName  string `json:"name" form:"name" valid:"required"`
	AAddr  string `json:"addr" form:"addr" valid:"required"`
	Weight int    `json:"weight" form:"weight" valid:"required"`
	Idx    string `json:"id" form:"-"`
}

func (s *formReverseSrver) ID() string      { return s.Idx }
func (s *formReverseSrver) SetID(id string) { s.Idx = id }
func (s *formReverseSrver) String() string {
	return fmt.Sprintf("formReverseSrver id: %s, Addr: %s", s.Idx, s.AAddr)
}
func (s formReverseSrver) Name() string { return s.NName }
func (s formReverseSrver) Addr() string { return s.AAddr }
func (s formReverseSrver) W() int       { return s.Weight }

type formCombReqCfg struct {
	SrvName string `json:"server_name" form:"server_name" valid:"required"` // http://ip:port/path?params
	PPath   string `json:"path" form:"path" valid:"required"`               // path `/request/path`
	FField  string `json:"field" form:"field" valid:"required"`             // want got field
	MMethod string `json:"method" form:"method" valid:"required"`           // want match method
	Idx     string `json:"id" form:"-"`
}

func (c *formCombReqCfg) ID() string      { return c.Idx }
func (c *formCombReqCfg) SetID(id string) { c.Idx = id }
func (c *formCombReqCfg) String() string {
	return fmt.Sprintf("formCombReqCfg id: %s, prefix: %s", c.Idx, c.PPath)
}
func (c *formCombReqCfg) ServerName() string { return c.SrvName }
func (c *formCombReqCfg) Path() string       { return c.PPath }
func (c *formCombReqCfg) Field() string      { return c.FField }
func (c *formCombReqCfg) Method() string     { return c.MMethod }

type formNocacher struct {
	Regexp string `json:"regular" form:"regular" valid:"required"`
	Idx    string `json:"id" form:"-"`
}

func (i *formNocacher) String() string  { return fmt.Sprintf("formNocacher: %s", i.Regexp) }
func (i *formNocacher) ID() string      { return i.Idx }
func (i *formNocacher) SetID(id string) { i.Idx = id }
func (i *formNocacher) Regular() string { return i.Regexp }
