package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type pluginsResp struct {
	code.CodeInfo
	Plugins []string `json:"plugins"`
}

// PluginsGET ...
func PluginsGET(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var (
		resp = new(pluginsResp)
	)
	resp.Plugins = []string{
		"plugin.proxy", "plugin.cache",
		"plugin.httplog", "plugin.ratelimit",
	}
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
