package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

// RolePOST ...
func RolePOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiRole)
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

	if err := Global().NewRole(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type rolesGetForm struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type rolesGetResp struct {
	code.CodeInfo
	Roles []*apiRole `json:"roles"`
	Total int        `json:"total"`
}

// RolesGET ...
func RolesGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(rolesGetForm)
		resp = new(rolesGetResp)
	)
	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}
	roles, total := Global().RolePage(form.Limit, form.Offset)
	resp.Roles = make([]*apiRole, len(roles))
	resp.Total = total
	for idx, role := range roles {
		resp.Roles[idx] = loadFormRole(role)
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// RolePUT ...
func RolePUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiRole)
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
	if err := Global().EditRole(id, form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// RoleDELETE ...
// func RoleDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		resp = new(commonResp)
// 	)
// 	id := param.ByName("id")
// 	if err := Global().DelRole(id); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }

// // RoleAssignPermPOST ...
// func RoleAssignPermPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		form = new(apiRole)
// 		resp = new(commonResp)
// 	)
// 	if err := Bind(form, req); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	if err := Valid(form); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	id := param.ByName("id")
// 	if err := Global().AssignPerm(id, form.PermIDs...); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }

// // RoleRevokePermPOST ...
// func RoleRevokePermPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		form = new(apiRole)
// 		resp = new(commonResp)
// 	)

// 	if err := Bind(form, req); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	if err := Valid(form); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	id := param.ByName("id")
// 	if err := Global().RevokePerm(id, form.PermIDs...); err != nil {
// 		ResponseWithError(w, resp, err)
// 		return
// 	}

// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }
