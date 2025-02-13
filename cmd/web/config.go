package main

import (
	"database/sql"
	"encoding/gob"
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
	_ = os.Getenv("addr")
	db, err := sql.Open("pgx", "host=localhost port=1234 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5")
	if err != nil {
		log.Panic("can not connect to db")
	}
	return db
}

func initSession() *scs.SessionManager {
	//allow user type to store in session
	gob.Register(data.User{})
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
			return redis.Dial("tcp", "127.0.0.1:6380")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return &redis
}
