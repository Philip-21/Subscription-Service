package main

import (
	"os"
	"os/signal"
	"syscall"
)

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
Shutdown() is responsible for performing any necessary
cleanup tasks and waiting for all running goroutines to complete
*/
func (app *Config) Shutdown() {
	//perform any cleanup task
	app.InfoLog.Println("would run cleaning up tasks....")

	//block until waitgroup is empty(counter hits 0)
	//waits until the goroutine executes
	app.Wait.Wait()

	app.Mailer.DoneChan <- true
	app.ErrorChanDone <- true

	app.InfoLog.Println("closing channels and shutting down application..")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
	close(app.ErrorChan)
	close(app.ErrorChanDone)
}

func (app *Config) ListenForErrors() {
	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorChanDone:
			return
		}
	}
}

// creating a mail server for testing
func (app *Config) createMail() Mail {
	//create channels
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100) //a buffered channel taking in 100messages before it locks
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromAddress: "philip@company.com",
		FromName:    "philip",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}
	return m
}
