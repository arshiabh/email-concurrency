package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/arshiabh/email-concurrency/cmd/web/data"
)

var pathToTemplate = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]string
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	User          *data.User
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	//parts should be render
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplate),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplate),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplate),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplate),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplate),
	}
	templateSlice := []string{}
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplate, t))
	templateSlice = append(templateSlice, partials...)

	if td == nil {
		td = &TemplateData{}
	}
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErroLogger.Fatal(err)
	}
	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErroLogger.Fatal(err)
	}
}

func (app *application) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true
		user, ok := app.Session.Get(r.Context(), "user").(data.User)
		if !ok {
			app.ErroLogger.Println("cannot get user from context")
		} else {
			td.User = &user
		}
	}
	td.Now = time.Now()
	return td
}

func (app *application) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
