package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"subscription-service/database"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

func initSession() *scs.SessionManager {
	gob.Register(database.User{}) //registers the struct to be used in the session
	// set up session
	session := scs.New()
	//store info for every session in reddis
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error in loading env file ", err)
	}
	conn := os.Getenv("REDIS")

	redisPool := &redis.Pool{
		MaxIdle: 10,
		//dial a redis server
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conn)
		},
	}

	return redisPool
}

// a secure session a used when a user is logged in
func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
