package main

import (
	"net/http"

	"github.com/arshiabh/email-concurrency/cmd/web/data"
)

func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *application) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	if err := r.ParseForm(); err != nil {
		app.ErroLogger.Println(err)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	var user *data.User
	user, err := app.Store.User.GetByEmail(email)
	if err != nil {
		app.InfoLogger.Println(err)
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//check it !! should call this method for user logged in
	valid, err := user.PasswordMatches(password)
	if err != nil {
		app.ErroLogger.Println(err)
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !valid {
		msg := Message{
			Data:    "invalid password",
			To:      user.Email,
			Subject: "invalid password",
		}
		app.sendEmail(msg)
		app.Session.Put(r.Context(), "error", "invalid password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "successfully logged in")
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *application) HandlePostRegister(w http.ResponseWriter, r *http.Request) {

}

func (app *application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) HandleActivateUser(w http.ResponseWriter, r *http.Request) {

}
