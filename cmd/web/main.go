package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscription-service/database"
	"sync"
	"syscall"
)

const webPort = "80"

func (app *Config) serve() {
	//start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server......")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

/*
gracefully shut down in response to certain signals and ensures
that all running goroutines complete their work before the application exits.
it closes channels and shuts down the app
*/

/*
ListenForShutdown() method sets up a channel called quit and registers a notification with the signal package
to listen for two signals: SIGINT and SIGTERM.
These signals indicate that the application should be shut down gracefully.
Once one of these signals is received,
the method blocks until the Shutdown() method is called
*/
func (app *Config) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	//Notify causes package signal to relay incoming signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Shutdown()

	os.Exit(0)
}

/*
Shutdown() method is responsible for performing any necessary
cleanup tasks and waiting for all running goroutines to complete
*/
func (app *Config) Shutdown() {
	//perform any cleanup task
	app.InfoLog.Println("would run cleaning up tasks....")

	//block until waitgroup is empty(counter hits 0)
	//waits until the goroutine executes
	app.Wait.Wait()

	app.InfoLog.Println("closing channels and shutting down application..")
}

func main() {

	//connect to db
	db := database.InitDB()

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
		Session:  session,
		DB:       db,
		Wait:     &wg,
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Models:   database.New(db),
	}

	//set up mail
	app.Mailer = app.createMail()
	go app.ListenForMail()

	//listen for signals
	go app.ListenForShutdown()

	//listen for web application
	app.serve()
}
