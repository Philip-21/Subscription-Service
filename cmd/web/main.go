package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"subscription-service/database"
	"sync"
	"time"
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
		app.ErrorLog.Panic("error in starting web server",err)
	}
	app.InfoLog.Println("Web Server Started")

}

func main() {

	//connect to db
	db := database.InitDB()
	//seeding actions
	err1 := database.Users(db)
	if err1 != nil {
		log.Printf("%s Couldnt seed Users table", err1)
		return
	}
	err2 := database.Plans(db)
	if err2 != nil {
		log.Printf("%s Couldnt seed plans table", err2)
		return
	}
	err3 := database.User_Plans(db)
	if err3 != nil {
		log.Printf("%s Couldnt Seed User-plans table ", err3)
		return
	}
	err4 := database.Altertable(db)
	if err4 != nil {
		log.Printf("%s couldnt alter user-plan table", err4)
		return
	}
	p := database.Plan{
		ID:         1,
		PlanName:   "Bronze Plan",
		PlanAmount: 5000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	m := database.Plan{
		ID:         2,
		PlanName:   "Silver Plan",
		PlanAmount: 20000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	t := database.Plan{
		ID:         3,
		PlanName:   "Gold Plan",
		PlanAmount: 100000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err1 = database.SeedUserPlans(db, &p)
	if err1 != nil {
		log.Printf("Error %s in Inserting Item in plan table", err1)
		return
	}
	err2 = database.SeedUserPlans(db, &m)
	if err2 != nil {
		log.Printf("Error %s in Inserting Item in Plan table", err2)
		return
	}
	err3 = database.SeedUserPlans(db, &t)
	if err3 != nil {
		log.Printf("Error %s in Inserting Item in Plan table", err3)
		return
	}

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
