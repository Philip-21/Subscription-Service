package main

import (
	"fmt"
	"html/template"
	"net/http"
	"subscription-service/database"
)

func (app *Config) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.go.html", nil)
}

func (app *Config) Loginpage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.go.html", nil)
}

func (app *Config) PostLoginpage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid Login ",
		}
		//call the channel func that sends email message easily
		app.sendemail(msg) //the app.Mailer.MailerChan which is a reciever from the message struct communictes this message sent to the
		//creating a session
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successful login!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.go.html", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	//validate user
	u := database.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		//specify the user is not an admin or not active
		Active:  0,
		IsAdmin: 0,
	}
	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create a user")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
	// send an activation email
	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedUrl := GenerateTokenFromString(url) //prevents url from being tampered
	app.InfoLog.Println(signedUrl)
	//create email message
	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account by verifying mail",
		Template: "confirmation-email",
		Data:     template.HTML(signedUrl), //cast the url into the html template
	}
	app.sendemail(msg)
	app.Session.Put(r.Context(), "flash", "Confirmation Mail Sent. Check your email to verify your mail.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// send activation email to confirm we have the valid email address
func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url
	url := r.RequestURI
	testUrl := fmt.Sprintf("http://localhost%s", url)
	okay := VerifyToken(testUrl) //the url with the hash appended to it
	if !okay {
		app.Session.Put(r.Context(), "error", "invalid token")
		app.ErrorLog.Println("invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//activate account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No User Found")
		app.ErrorLog.Println("No User Found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable User to update User.")
		app.ErrorLog.Println("Unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "flash", "Account activated . Ypu can login.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached

	//subscribe the User is an Account
}
func (app *Config) SubcribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is chosen

	// get the plan from the database

	// get the user from the session

	// generate an invoice

	// send an email with the invoice attached

	// generate a manual

	// send an email with the manual attached

	// subscribe the user to an account

	// redirect
}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	if !app.Session.Exists(r.Context(), "userID") {
		app.Session.Put(r.Context(), "warning", "Login to access this page")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	//pass the var to the template
	dataMap := make(map[string]any)
	dataMap["plans"] = plans
	app.render(w, r, "plans.page.go.html", &TemplateData{
		Data: dataMap,
	})
}
