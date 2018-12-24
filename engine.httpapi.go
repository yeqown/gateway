package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yeqown/gateway/config/api"
	"github.com/yeqown/gateway/plugin"
	"github.com/yeqown/gateway/utils"
	"github.com/yeqown/server-common/code"
)

type httpapiPlg struct {
	Name    string           `json:"name"`
	Enabled bool             `json:"enabled"`
	Status  plugin.PlgStatus `json:"status"`
	Idx     int              `json:"idx"`
}

type pluginsGetResp struct {
	code.CodeInfo
	AllPlugins []httpapiPlg `json:"plugins"`
	Num        int          `json:"num"`
}

func (e *Engine) pluginsGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(pluginsGetResp)
	)

	resp.Num = e.numAllPlugin
	resp.AllPlugins = make([]httpapiPlg, e.numAllPlugin)

	for idx, plg := range e.allPlugins {
		resp.AllPlugins[idx] = httpapiPlg{
			Name:    plg.Name(),
			Enabled: plg.Enabled(),
			Status:  plg.Status(),
			Idx:     idx,
		}
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}

type enablePlgForm struct {
	Enabled bool `form:"enabled"`
	Idx     int  `form:"idx"` // id means plugin position in active plugins, 0 means the first one.
}

type enablePlgResp struct {
	code.CodeInfo
}

// TODO: maybe lock the engine for change active plugins list
func (e *Engine) enablePluginGET(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	var (
		resp = new(enablePlgResp)
		form = new(enablePlgForm)
	)

	if err := api.Bind(form, req); err != nil {
		api.ResponseWithError(w, resp, err)
		return
	}

	if err := api.Valid(form); err != nil {
		api.ResponseWithError(w, resp, err)
		return
	}

	e.enablePlugin(form.Enabled, form.Idx)

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	utils.ResponseJSON(w, resp)
}
