package application


import (
	"fmt"
	"context"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	// create server
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}
	// handle errors
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server : %w", err)
	}

	return nil
}
