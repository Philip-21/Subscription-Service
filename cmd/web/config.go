package main

import (
	"log"
	"subscription-service/database"
	"sync"

	"github.com/alexedwards/scs/v2"
	"gorm.io/gorm"
)

type Config struct {
	Session       *scs.SessionManager
	DB            *gorm.DB
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Wait          *sync.WaitGroup
	Models        database.Models
	Mailer        Mail //defining the mail struct for the mail server
	ErrorChan     chan error
	ErrorChanDone chan bool
}
