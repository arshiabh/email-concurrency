package data

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Store struct {
	User User
	Plan Plan
}

var db *sql.DB

func New(dbPool *sql.DB) Store {
	db = dbPool

	return Store{
		User: User{},
		Plan: Plan{},
	}
}
