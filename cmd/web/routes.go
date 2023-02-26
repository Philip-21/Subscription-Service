package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	//create a router
	mux := chi.NewRouter()

	//middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	//app routes
	mux.Get("/", app.homePage)

	mux.Get("/login", app.Loginpage)
	mux.Post("/login", app.PostLoginpage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate", app.ActivateAccount)

	//Mount attaches another http.Handler or chi Router as a subrouter along a routing
	//attaches handlers for only authenticated users
	mux.Mount("/members", app.authRouter())

	//Sending Email Synchronously (not in use just for reference)
	mux.Get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "philip@company.com",
			FromName:    "philip",
			//make an error chan
			ErrorChan: make(chan error),
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "testing mail service",
			Data:    "Hello World",
		}
		//sending email
		m.SendMail(msg, make(chan error))
	})

	return mux
}

func (app *Config) authRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.Auth)

	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubcribeToPlan)
	return mux

}
