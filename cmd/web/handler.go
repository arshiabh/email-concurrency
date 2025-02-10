package main

import "net/http"

func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *application) HandlePostLogin(w http.ResponseWriter, r *http.Request) {

}

func (app *application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *application) HandlePostRegister(w http.ResponseWriter, r *http.Request) {

}

func (app *application) HandleLogout(w http.ResponseWriter, r *http.Request) {

}

func (app *application) HandleActivateUser(w http.ResponseWriter, r *http.Request) {

}
