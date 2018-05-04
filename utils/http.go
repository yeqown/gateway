package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// CopyRequest copy request from a http.Request
func CopyRequest(req *http.Request) *http.Request {
	body, _ := ioutil.ReadAll(req.Body)
	rdOnly := ioutil.NopCloser(bytes.NewBuffer(body))

	reqCpy, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	reqCpy.Header = req.Header
	req.Body = rdOnly
	return reqCpy
}

// ParseRequestForm ...
// parse request and get form form body or url
func ParseRequestForm(cpyReq *http.Request) map[string]interface{} {
	reqData := make(map[string]interface{})
	switch cpyReq.Method {
	case http.MethodPost, http.MethodPut:
		cpyReq.ParseMultipartForm(32 << 20)
	case http.MethodGet:
		cpyReq.ParseForm()
	default:
		cpyReq.ParseForm()
	}
	for k, v := range cpyReq.Form {
		reqData[k] = v
	}
	return reqData
}
