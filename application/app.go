package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	psql   *sql.DB
}

func New() *App {
	err := godotenv.Load()

	if err != nil {
		println(".env doesnt exist!", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		println("error connecting to database", err)
	}
	// driver, err := postgres.WithInstance(db, &postgres.Config{})

	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
		psql:   db,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	// create server
	server := &http.Server{
		Addr:    ":3000",
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
