package main

import "net/http"

func (app *Config) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) Loginpage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginpage(w http.ResponseWriter, r *http.Request) {}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {}

//send activation email to confirm we have the valid email address
func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {}
