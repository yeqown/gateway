package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type permsGetForm struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type permsGetResp struct {
	code.CodeInfo
	Perms []*apiPermission `json:"permissions"`
	Total int              `json:"total"`
}

// PermissionsGET ...
func PermissionsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(permsGetForm)
		resp = new(permsGetResp)
	)

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	perms, total := Global().PermissionPage(form.Limit, form.Offset)
	resp.Perms = make([]*apiPermission, len(perms))
	resp.Total = total
	for idx, perm := range perms {
		resp.Perms[idx] = loadFromPermission(perm)
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// PermissionDELETE ...
func PermissionDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	id := param.ByName("id")
	if err := Global().DelPermission(id); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// PermissionPUT ...
func PermissionPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPermission)
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
	id := param.ByName("id")
	if err := Global().EditPermission(id, form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// PermissionPOST ...
func PermissionPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPermission)
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

	if err := Global().NewPermission(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
