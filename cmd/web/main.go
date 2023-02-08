package main

import (
	"log"
	"os"
	"sync"
)

const webPort = "80"

func main() {

	//connect to db
	db := initDB()
	//db.Ping()

	//create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//createb sessions to render templates
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
	}

	//set up mail

	//listen for web application
}
