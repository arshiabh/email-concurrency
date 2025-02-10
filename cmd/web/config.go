package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/arshiabh/email-concurrency/cmd/web/data"
	"github.com/gomodule/redigo/redis"
)

type application struct {
	Session    *scs.SessionManager
	DB         *sql.DB
	InfoLogger *log.Logger
	ErroLogger *log.Logger
	Wait       *sync.WaitGroup
	Store      *data.Store
}

func initDB() *sql.DB {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Panic("can not connect to db")
	}
	return db
}

func initSession() *scs.SessionManager {
	session := scs.New()
	//store all information from session in redis
	session.Store = redisstore.New(initRedis())
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session
}

func initRedis() *redis.Pool {
	redis := redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return &redis
}
