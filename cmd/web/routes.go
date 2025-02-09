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
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", app.Index)
	})
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
