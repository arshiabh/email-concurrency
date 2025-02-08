package main

import (
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db := initDB()
	session := initSession()
	var wg sync.WaitGroup

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		Session:    session,
		DB:         db,
		Wait:       &wg,
		InfoLogger: infoLogger,
		ErroLogger: errorLogger,
	}
	mux := app.mount()
	if err := app.run(mux); err != nil {
		app.ErroLogger.Fatal(err)
	}
}
