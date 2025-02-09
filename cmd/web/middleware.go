package main

import "net/http"

func (app *application) sessionload(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}
