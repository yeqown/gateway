package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type proxycfgReverseSrvsForm struct {
	Limit  int `form:"limit" valid:"gte=0,lte=10"`
	Offset int `form:"offset" valid:"gte=0"`
}

type proxycfgReverseSrvsGroup struct {
	code.CodeInfo
	Group []rule.ReverseServer `json:"group"`
	Total int                  `json:"total"`
}

// ProxyConfigReverseSrvGroupGET ...
func ProxyConfigReverseSrvGroupGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(proxycfgReverseSrvsForm)
		resp = new(proxycfgReverseSrvsGroup)
	)

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if form.Limit == 0 {
		form.Limit = 10
	}

	group := param.ByName("group")
	resp.Group = Global().ReverseServerGroup(group, form.Offset, form.Limit+form.Offset)
	resp.Total = Global().ReverseServerGroupPageCount(group)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvGroupDELETE ...
func ProxyConfigReverseSrvGroupDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)
	group := param.ByName("group")
	if err := Global().DelReverseServerGroup(group); err != nil {
		responseWithError(w, resp, err)
		return
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type proxycfgReverseSrvResp struct {
	code.CodeInfo
	Rule rule.ReverseServer `json:"rule"`
}

// ProxyConfigReverseSrvGET ...
func ProxyConfigReverseSrvGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(proxycfgReverseSrvResp)
	)

	group, id := param.ByName("group"), param.ByName("id")
	resp.Rule = Global().ReverseServerByID(group, id)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvPOST ... [fixed]TOFIX: filestore 不展示ID
func ProxyConfigReverseSrvPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(formReverseSrver)
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

	group := param.ByName("group")
	if err := Global().NewReverseServer(group, form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvPUT ...
func ProxyConfigReverseSrvPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(formReverseSrver)
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

	group, idxs := param.ByName("group"), param.ByName("id")
	id := fmt.Sprintf("%s#%s", group, idxs)
	if err := Global().UpdateReverseServer(id, form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvDELETE ...
func ProxyConfigReverseSrvDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	group, idxs := param.ByName("group"), param.ByName("id")
	id := fmt.Sprintf("%s#%s", group, idxs)

	if err := Global().DelReverseServer(id); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
