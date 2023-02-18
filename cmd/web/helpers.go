package main

import (
	"os"
	"os/signal"
	"syscall"
)

// a  helper wrapper to send email easily
func (app *Config) sendemail(msg Message) {
	//add counter to waitgroup , increment wg by 1
	app.Wait.Add(1)
	app.Mailer.MailerChan <- msg //send message to the mail channel
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

	app.Mailer.DoneChan <- true

	app.InfoLog.Println("closing channels and shutting down application..")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
}
