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
	Cfg *presistence.Instance `json:"config"`
}

// AllConfigsGET get all configs
func AllConfigsGET(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var (
		resp = new(configsResp)
	)

	resp.Cfg = Global().Instance()
	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
	return
}

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

type gateConfigResp struct {
	code.CodeInfo
	Logpath string `json:"logpath"`
	Port    int    `json:"port"`
}

// GateConfigGET ...
func GateConfigGET(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var (
		resp = new(gateConfigResp)
	)

	resp.Logpath = Global().Instance().Logpath
	resp.Port = Global().Instance().Port

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type updateGateConfigResp struct {
	Logpath string `form:"logpath" valid:"required"`
	Port    int    `form:"port" valid:"required"`
}

// GateConfigPUT ...
func GateConfigPUT(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var (
		resp = new(commonResp)
		form = new(updateGateConfigResp)
	)

	if err := bind(form, req); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := valid(form); err != nil {
		responseWithError(w, resp, err)
		return
	}

	if err := Global().UpdateGateConfig(form.Logpath, form.Port); err != nil {
		responseWithError(w, resp, err)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
