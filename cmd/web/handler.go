package main

import "net/http"

func (app *application) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
