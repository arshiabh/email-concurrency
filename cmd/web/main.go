package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/arshiabh/email-concurrency/cmd/web/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db := initDB()
	session := initSession()
	//use the same waitgroup for rest service to avoid error
	var wg sync.WaitGroup

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	store := data.New(db)
	mail := createMail(&wg)

	app := &application{
		Session:    session,
		DB:         db,
		Wait:       &wg,
		InfoLogger: infoLogger,
		ErroLogger: errorLogger,
		Store:      &store,
		Mailer:     mail,
	}
	//listern for shutdown
	go app.listenForShutdown()
	//listen for email
	go app.listernForEmail()
	mux := app.mount()
	if err := app.run(mux); err != nil {
		app.ErroLogger.Fatal(err)
	}
}

func (app *application) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	//block and wait to recive
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *application) shutdown() {
	app.InfoLogger.Println("cleaning for shutdown")
	// wait for other go routine
	app.Wait.Wait()
	//after wait all wg done trigger done to clean
	app.Mailer.DoneChan <- true

	app.InfoLogger.Println("closing server")
	//close chan after to avoid getting data
	close(app.Mailer.MailerChan)
	close(app.Mailer.DoneChan)
	close(app.Mailer.ErrorChan)
}
