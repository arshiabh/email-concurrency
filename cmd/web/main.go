package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)


func main() {
	db := initDB()
	db.Ping()
}

func initDB() *sql.DB {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Panic("can not connect to db")
	}
	return db
}
