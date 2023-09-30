package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mypassword"
	dbname   = "postgres"
)

func New() *App {
	// load env variables
	err := godotenv.Load()
	// check if .env exists
	if err != nil {
		println(".env doesnt exist!", err)
	}
	// connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// connect
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		println("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		println("error pinging db")
	}
	defer db.Close()
	// create app instance
	app := &App{
		db: db,
	}
	// load routes
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	// create server
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	// handle SERVER errors
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server : %w", err)
	}

	return nil
}
