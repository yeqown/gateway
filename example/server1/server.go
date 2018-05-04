package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type response struct {
	Msg  string `json:"data"`
	Code int    `json:"code"`
}

func main() {
	http.HandleFunc("/srv/name", func(w http.ResponseWriter, req *http.Request) {
		req.ParseMultipartForm(32 << 10)
		name := req.Form.Get("name")
		msg := name
		if name == "" {
			msg = "error, empty name value"
		}
		bs, _ := json.Marshal(response{Msg: msg, Code: 0})
		fmt.Fprintf(w, string(bs))
		return
	})

	http.HandleFunc("/srv/id", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		id := req.Form.Get("id")
		msg := id
		if id == "" {
			msg = "error, empty id"
		}

		bs, _ := json.Marshal(response{Msg: msg, Code: 0})
		fmt.Fprintf(w, string(bs))
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		bs, _ := json.Marshal(response{Msg: "this is server 1", Code: 0})
		fmt.Fprintf(w, string(bs))
		return
	})

	log.Printf("listen on: %s\n", ":8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
