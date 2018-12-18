package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type configsResp struct {
	code.CodeInfo
	cfg *presistence.Instance
}

// AllConfigsGET get all configs
func AllConfigsGET(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var (
		resp = new(configsResp)
	)

	resp.cfg = Global().Instance()
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
	return
}
