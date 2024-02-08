package main

import (
	"database/sql"
	"encoding/gob"
	"final_project/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8085"

func main() {
	// connect to the database
	db := initDB()

	// create sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(db),
	}

	// set up mail
	app.Mailer = app.createMail()
	go app.listenForMail()

	// listen for signals
	go app.listenForShutdown()

	// listen for web connections
	app.serve()
}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Printf("Starting web server at %s...", webPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}

	return conn
}

func connectToDB() *sql.DB {
	count := 0
	dsn := getOsEnvOrDefault("DSN", "host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready ...")
		} else {
			log.Println("connected to database")
			return connection
		}

		count++

		if count > 10 {
			return nil
		}

		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil

}

func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	// set up session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session

}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", getOsEnvOrDefault("REDIS", "127.0.0.1:6379"))
		},
	}

	return redisPool
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// Run all cleanup tasks
	app.InfoLog.Println("Would run cleanup tasks...")

	// block until waitgroup empty
	app.Wait.Wait()

	// send done channel as true for all mails
	app.Mailer.DoneChan <- true
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)

	// Close now
	app.InfoLog.Println("closing channels and shutting down application")
}

func (app *Config) createMail() Mail {
	// create channe;s
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}
	return m
}

func getOsEnvOrDefault(key string, fallback string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return fallback
	}
	return val
}