package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

// UserPOST ...
func UserPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiUser)
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

	if err := Global().NewUser(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type usersGetForm struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type usersGetResp struct {
	code.CodeInfo
	Users []*apiUser `json:"users"`
	Total int        `json:"total"`
}

// UsersGET ...
func UsersGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(usersGetForm)
		resp = new(usersGetResp)
	)
	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}
	users, total := Global().UserPage(form.Limit, form.Offset)
	resp.Users = make([]*apiUser, len(users))
	resp.Total = total
	for idx, user := range users {
		resp.Users[idx] = loadFormUser(user)
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// UserPUT ...
func UserPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiUser)
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
	if err := Global().EditUser(id, form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// UserDELETE ...
func UserDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	id := param.ByName("id")
	if err := Global().DelUser(id); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// UserAssignRolePOST ...
// func UserAssignRolePOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		form = new(apiUser)
// 		resp = new(commonResp)
// 	)

// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }

// // UserRevokeRolePOST ...
// func UserRevokeRolePOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
// 	var (
// 		form = new(apiUser)
// 		resp = new(commonResp)
// 	)

// 	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
// 	utils.ResponseJSON(w, resp)
// }
