package main

import (
	"fmt"
	"html/template"
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
	err := r.ParseForm()
	if err != nil {
		app.ErroLogger.Println(err)
	}
	//to-do validate form value we get
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		IsAdmin:   0,
		Active:    0,
	}

	_, err = app.Store.User.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "cannot create user!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedurl := GenerateTokenFromString(url)
	app.InfoLogger.Println(signedurl)

	msg := Message{
		To:       u.Email,
		Subject:  "User Activation",
		Template: "confirmation-email",
		Data:     template.HTML(signedurl),
	}

	app.sendEmail(msg)
	app.Session.Put(r.Context(), "flash", "check your email for activation")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) HandleActivateUser(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	app.InfoLogger.Println(url)
	testuri := fmt.Sprintf("http://localhost%s", url)

	valid := VerifyToken(testuri)

	if !valid {
		app.Session.Put(r.Context(), "error", "invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user, err := app.Store.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "user not found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user.Active = 1
	if err := user.Update(); err != nil {
		app.Session.Put(r.Context(), "error", "user cannot updated")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "user activated. now you can login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
