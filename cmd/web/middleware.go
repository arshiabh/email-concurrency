package main

import "net/http"

func (app *application) sessionload(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			app.Session.Put(r.Context(), "error", "first login!")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
