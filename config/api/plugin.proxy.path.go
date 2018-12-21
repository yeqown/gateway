package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

var (
	errJSONDataNeed = errors.New("need JSON body data")
)

type proxyPathsGetForm struct {
	Limit  int `form:"limit" valid:"gte=0,lte=10"`
	Offset int `form:"offset" valid:"gte=0"`
}

type proxyPathsGetResp struct {
	code.CodeInfo
	Total int             `json:"total"`
	Rules []*apiPathRuler `json:"rules"`
}

// ProxyConfigPathsGET ...
func ProxyConfigPathsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(proxyPathsGetResp)
		form = new(proxyPathsGetForm)
	)

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	rules, total := Global().PathRulesPage(form.Offset, form.Limit)
	resp.Total = total
	for _, r := range rules {
		resp.Rules = append(resp.Rules, loadFromPathRuler(r))
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type proxyPathGetResp struct {
	code.CodeInfo
	Rule *apiPathRuler `json:"rule,omitempty"`
}

// ProxyConfigPathGET ...
func ProxyConfigPathGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(proxyPathGetResp)
	)
	id := param.ByName("id")
	rule := Global().PathRuleByID(id)
	resp.Rule = loadFromPathRuler(rule)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigPathPOST ... [fixed] TOFIX request by JSON
func ProxyConfigPathPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPathRuler)
		resp = new(commonResp)
	)

	byts, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseWithError(w, resp, errJSONDataNeed)
		return
	}
	err = json.Unmarshal(byts, form)
	if err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := Global().NewPathRule(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigPathPUT ... [fixed] TOFIX request by JSON
func ProxyConfigPathPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPathRuler)
		resp = new(commonResp)
	)

	byts, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseWithError(w, resp, errJSONDataNeed)
		return
	}
	err = json.Unmarshal(byts, form)
	if err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}
	id := param.ByName("id")
	if err := Global().UpdatePathRule(id, form); err != nil {
		responseWithError(w, resp, err)
		return
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigPathDELETE ...
func ProxyConfigPathDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	id := param.ByName("id")

	if err := Global().DelPathRule(id); err != nil {
		responseWithError(w, resp, err)
		return
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
