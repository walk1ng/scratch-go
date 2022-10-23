package main

import "net/http"

func hellohandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
