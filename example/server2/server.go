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
		bs, _ := json.Marshal(response{Msg: "server2", Code: 0})
		fmt.Fprintf(w, string(bs))
		return
	})

	http.HandleFunc("/srv/id", func(w http.ResponseWriter, req *http.Request) {
		bs, _ := json.Marshal(response{Msg: "2", Code: 0})
		fmt.Fprintf(w, string(bs))
		return
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	bs, _ := json.Marshal(response{Msg: "this is server 2", Code: 0})
	// 	fmt.Fprintf(w, string(bs))
	// 	return
	// })

	log.Printf("listen on: %s\n", ":8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
