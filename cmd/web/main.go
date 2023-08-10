package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"subscription-service/database"
	"sync"
)

const webPort = "8080"

func (app *Config) serve() {
	//start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server......")
	err := srv.ListenAndServe()
	if err != nil {
		app.ErrorLog.Panic("error in starting web server", err)
	}
	app.InfoLog.Println("Web Server Started")

}

func main() {

	//connect to db
	db := database.ConnectToDB()
	//seeding actions

	//create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//create sessions to render templates
	session := initSession()

	//create channels

	//create waitgroups
	wg := sync.WaitGroup{}

	//set up application configurations
	app := Config{
		Session:       session,
		DB:            db,
		Wait:          &wg,
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Models:        database.New(db),
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	//set up mail
	app.Mailer = app.createMail()
	//listens and sends mail asynchronously
	go app.ListenForMail()

	//listen for signals
	go app.ListenForShutdown()

	//listen for errors
	go app.ListenForErrors()

	//listen for web application
	app.serve()

}
