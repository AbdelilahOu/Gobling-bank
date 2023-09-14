package application


import (
	"fmt"
	"context"
	"net/http"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{})
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	// create server
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}
	// redis health check
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("error connecting to redis : %w", err)
	}
	// handle SERVER errors
	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server : %w", err)
	}

	return nil
}
