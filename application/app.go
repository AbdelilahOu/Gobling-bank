package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func New() *App {
	// load env variables
	err := godotenv.Load()
	// check if .env exists
	if err != nil {
		println(".env doesnt exist!", err)
	}
	//
	connect, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		println("error connecting to db")
	}
	// create app instance
	app := &App{
		db: connect,
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
