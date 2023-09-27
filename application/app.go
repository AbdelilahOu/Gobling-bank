package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

type App struct {
	router http.Handler
}

func New() *App {
	// load env variables
	err := godotenv.Load()
	// check if .env exists
	if err != nil {
		println(".env doesnt exist!", err)
	}
	// create app instance
	app := &App{
		router: loadRoutes(),
	}

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
