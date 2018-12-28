package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

// PermitURLPOST ...
func PermitURLPOST(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPermitURL)
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
	if err := Global().NewPermitURL(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type permitURLGetForm struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type permitURLGetResp struct {
	code.CodeInfo
	PermitURLs []*apiPermitURL `json:"permit_urls"`
	Total      int             `json:"total"`
}

// PermitURLsGET ...
func PermitURLsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(permitURLGetForm)
		resp = new(permitURLGetResp)
	)

	if err := Bind(form, req); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	if err := Valid(form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	urls, total := Global().PermitURLPage(form.Limit, form.Offset)
	resp.Total = total
	resp.PermitURLs = make([]*apiPermitURL, len(urls))
	for idx, url := range urls {
		resp.PermitURLs[idx] = loadFormPermitURL(url)
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// PermitURLPUT ...
func PermitURLPUT(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		form = new(apiPermitURL)
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
	if err := Global().EditPermitURL(id, form); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

// PermitURLDELETE ...
func PermitURLDELETE(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(commonResp)
	)

	id := param.ByName("id")
	if err := Global().DelPermitURL(id); err != nil {
		ResponseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
