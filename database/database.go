package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	counts := 0
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dsn := os.Getenv("DSN")
	fmt.Println(dsn)
	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Print("connected to database!")
			return connection
		}
		if counts > 10 {
			return nil
		}
		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{}, &Plan{})
	if err != nil {
		log.Println("Error in Migrations")
	}
	log.Println("migrations successful")
	//verifies if a connection to the database is still alive, establishing a connection if necessary.
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("error in getting sql")
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Println("Error in connection", err)
		return nil, err
	}
	return db, nil
}
