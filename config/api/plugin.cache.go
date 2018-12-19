package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type cacheConfigsForm struct {
	Limit  int `form:"limit" valid:"gte=0,lte=10"`
	Offset int `form:"offser" valid:"gte=0"`
}

type cacheConfigsResp struct {
	code.CodeInfo
	Rules []rule.Nocacher `json:"rule,omitempty"`
	Total int             `json:"total,omitempty"`
}

// CacheConfigsGET ...
func CacheConfigsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(cacheConfigsForm)
		resp = new(cacheConfigsResp)
	)

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	resp.Rules = Global().NocacheRules(form.Offset, form.Offset+form.Limit)
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// single
type cacheConfigResp struct {
	code.CodeInfo
	Rule rule.Nocacher `json:"rule,omitempty"`
}

// CacheConfigGET ...
func CacheConfigGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(cacheConfigResp)
	)

	id := param.ByName("id")
	resp.Rule = Global().NocacheRuleByID(id)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// CacheConfigPOST ...
func CacheConfigPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(formNocacher)
		resp = new(commonResp)
	)

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := Global().NewNocacheRule(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// CacheConfigPUT ...
func CacheConfigPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(formNocacher)
		resp = new(commonResp)
	)
	id := param.ByName("id")

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := Global().UpdateNocacheRule(id, form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// CacheConfigDELETE ...
func CacheConfigDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	id := param.ByName("id")
	if err := Global().DelNocacheRule(id); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
