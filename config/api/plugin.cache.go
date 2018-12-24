package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type cacheConfigsForm struct {
	Limit  int `form:"limit" valid:"gte=0,lte=10"`
	Offset int `form:"offset" valid:"gte=0"`
}

type cacheConfigsResp struct {
	code.CodeInfo
	Rules []*apiNocacher `json:"rules"`
	Total int            `json:"total"`
}

// CacheConfigsGET ...
func CacheConfigsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(cacheConfigsForm)
		resp = new(cacheConfigsResp)
	)

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	rules, total := Global().NocacheRules(form.Offset, form.Limit)
	resp.Total = total
	for _, r := range rules {
		resp.Rules = append(resp.Rules, loadFromNocacher(r))
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// single
type cacheConfigResp struct {
	code.CodeInfo
	Rule *apiNocacher `json:"rule,omitempty"`
}

// CacheConfigGET ...
func CacheConfigGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(cacheConfigResp)
	)

	id := param.ByName("id")
	rule := Global().NocacheRuleByID(id)
	resp.Rule = loadFromNocacher(rule)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// CacheConfigPOST ...
func CacheConfigPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiNocacher)
		resp = new(commonResp)
	)

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Global().NewNocacheRule(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// CacheConfigPUT ...
func CacheConfigPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiNocacher)
		resp = new(commonResp)
	)
	id := param.ByName("id")

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Global().UpdateNocacheRule(id, form); err != nil {
		ResponseWithError(w, resp, err)
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
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
