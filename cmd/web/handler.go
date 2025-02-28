package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/arshiabh/email-concurrency/cmd/web/data"
	"github.com/phpdave11/gofpdf"
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

func (app *application) HandleChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Store.Plan.GetAll()
	if err != nil {
		app.ErroLogger.Println(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans
	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *application) HandleSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	user, ok := app.Session.Get(r.Context(), "userID").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "first login!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id")
	planID, _ := strconv.Atoi(id)
	plan, err := app.Store.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "unable to find plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		invoice, err := app.invoice(user, plan)
		if err != nil {
			app.ErrorChan <- err
		}
		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice",
			Data:     invoice,
			Template: "invoice",
		}
		app.sendEmail(msg)
	}()

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		pdf := app.generatePDF()

	}()

}

func (app *application) generatePDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)
	return pdf
}

func (app *application) invoice(u data.User, plan *data.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}
