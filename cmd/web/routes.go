package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(app.sessionload)
	r.Group(func(r chi.Router) {
		r.Get("/", app.HandleHome)
		r.Get("/login", app.HandleLogin)
		r.Post("/login", app.HandlePostLogin)
		r.Get("/register", app.HandleRegister)
		r.Post("/register", app.HandlePostRegister)
		r.Get("/logout", app.HandleLogout)
		r.Get("/activate", app.HandleActivateUser)
		//takes another handle to it self
		r.Mount("/members", app.authRouter())
	})
	return r
}

func (app *application) authRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(app.Auth)
	r.Get("/subscribe", app.HandleSubscribeToPlan)
	r.Get("/plans", app.HandleChooseSubscription)
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return srv.ListenAndServe()
}
