package main

import (
    "flag"
    "net/http"
)

func main() {

        flag.Parse()

	http.HandleFunc("/v1/authorize", Authorize)
	http.HandleFunc("/v1/sign", Sign)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func Authorize(w http.ResponseWriter, req *http.Request) {

}

func Sign(w http.ResponseWriter, req *http.Request) {

}
