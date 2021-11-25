package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type AppStatus struct {
	Status     string `json:"status"`
	Enviroment string `json:"enviroment"`
	Version    string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var configuration config
	flag.IntVar(&configuration.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&configuration.env, "env", "development", "Application enviroment (development | production")

	flag.StringVar(&configuration.db.dsn, "dsn", "postgres://postgres:admin@localhost/go_movies?sslmode=disable", "Postgress connection")

	db, err := openDB(configuration)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: configuration,
		logger: logger,
		models: models.NewModel(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", configuration.port)

	err = srv.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}

func openDB(configuration config) (*sql.DB, error) {
	db, err := sql.Open("postgres", configuration.db.dsn)

	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(context)

	if err != nil {
		return nil, err
	}

	return db, nil
}
