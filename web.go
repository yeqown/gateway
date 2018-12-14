package gateway

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"
// 	"path"
// 	"strings"

// 	"github.com/yeqown/gateway/logger"
// )

// // Prefix of URL
// const webPrefix = "/html/"

// // HTMLSrv ...
// var HTMLSrv *HTMLServer

// func init() {
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", IndexPage)
// 	mux.HandleFunc("/p1", func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintf(w, "page here!")
// 		return
// 	})
// 	HTMLSrv = &HTMLServer{
// 		mux:      mux,
// 		HTMLRoot: "../../web/html",
// 	}
// }

// // HTMLServer ...
// type HTMLServer struct {
// 	mux      *http.ServeMux
// 	HTMLRoot string
// }

// // ServeHTTP ...
// func (s *HTMLServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	logger.Logger.Info(req.URL.Path)

// 	req.URL.Path = strings.TrimPrefix(req.URL.Path, webPrefix)
// 	req.URL.Path = "/" + req.URL.Path
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	s.mux.ServeHTTP(w, req)
// }

// // IndexPage ...
// func IndexPage(w http.ResponseWriter, req *http.Request) {
// 	fp := path.Join(HTMLSrv.HTMLRoot, "index.html")
// 	tmpl := template.Must(template.ParseFiles(fp))
// 	tmpl.Execute(w, nil)
// 	return
// }
