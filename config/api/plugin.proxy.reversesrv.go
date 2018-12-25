package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type proxycfgReverseSrvGroupsResp struct {
	code.CodeInfo
	Groups map[string]int `json:"groups"`
}

// ProxyConfigReverseSrvGroupsGET ...
func ProxyConfigReverseSrvGroupsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(proxycfgReverseSrvGroupsResp)
	)
	resp.Groups = Global().ReverseServerGroups()

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type proxycfgReverseSrvsForm struct {
	Limit  int `form:"limit" valid:"gte=0,lte=10"`
	Offset int `form:"offset" valid:"gte=0"`
}

type proxycfgReverseSrvsGroup struct {
	code.CodeInfo
	Group []*apiReverseSrver `json:"group"`
	Total int                `json:"total"`
}

// ProxyConfigReverseSrvGroupGET ...
func ProxyConfigReverseSrvGroupGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(proxycfgReverseSrvsForm)
		resp = new(proxycfgReverseSrvsGroup)
	)

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	group := param.ByName("group")
	rules, total := Global().ReverseServerByGroup(group, form.Offset, form.Limit)
	resp.Total = total
	for _, r := range rules {
		resp.Group = append(resp.Group, loadFromReverseServer(r))
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type proxycfgReverseSrvGroupPutForm struct {
	Newname string `form:"newname" valid:"required"`
}

// ProxyConfigReverseSrvGroupPUT ...
func ProxyConfigReverseSrvGroupPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(proxycfgReverseSrvGroupPutForm)
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
	group := param.ByName("group")
	if err := Global().UpdateReverseServerGroupName(group, form.Newname); err != nil {
		ResponseWithError(w, resp, err)
		return
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvGroupPOST ...
// func ProxyConfigReverseSrvGroupPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		resp = new(commonResp)
// 	)
// 	group := param.ByName("group")
// 	if err := Global().NewReverseServerGroup(group); err != nil {
// 		responseWithError(w, resp, err)
// 		return
// 	}
// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }

// ProxyConfigReverseSrvGroupDELETE ...
func ProxyConfigReverseSrvGroupDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)
	group := param.ByName("group")
	if err := Global().DelReverseServerGroup(group); err != nil {
		ResponseWithError(w, resp, err)
		return
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type proxycfgReverseSrvResp struct {
	code.CodeInfo
	Rule *apiReverseSrver `json:"rule"`
}

// ProxyConfigReverseSrvGET ...
func ProxyConfigReverseSrvGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(proxycfgReverseSrvResp)
	)

	group, id := param.ByName("group"), param.ByName("id")
	rule := Global().ReverseServerByID(group, id)
	resp.Rule = loadFromReverseServer(rule)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvPOST ... 新增反向代理配置
func ProxyConfigReverseSrvPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiReverseSrver)
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

	group := param.ByName("group")
	if err := Global().NewReverseServer(group, form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// ProxyConfigReverseSrvPUT ...
func ProxyConfigReverseSrvPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiReverseSrver)
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

	_, id := param.ByName("group"), param.ByName("id")
	// id := fmt.Sprintf("%s#%s", group, idxs)
	if err := Global().UpdateReverseServer(id, form); err != nil {
		ResponseWithError(w, resp, err)
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

	_, id := param.ByName("group"), param.ByName("id")
	// id := fmt.Sprintf("%s#%s", group, idxs)

	if err := Global().DelReverseServer(id); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
