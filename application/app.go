package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
)

type App struct {
	router http.Handler
	psql   *sql.DB
}

func New() *App {
	// load env variables
	err := godotenv.Load()
	// check if .env exists
	if err != nil {
		println(".env doesnt exist!", err)
	}
	// created db connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// check if theres errro
	if err != nil {
		println("error connecting to database", err.Error())
	}
	// create app instance
	app := &App{
		router: loadRoutes(),
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
	// health check the db connection
	err := a.psql.Ping()
	if err != nil {
		println("error pinging db", err)
	}
	// created driver & check for error
	driver, err := postgres.WithInstance(a.psql, &postgres.Config{})
	if err != nil {
		println("error creating driver", err)
	}
	// create migration instance & check if theres error
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)
	if err != nil {
		println("error creating migration instance")
	}
	// run migration
	if err := m.Up(); err != nil {
		println("error running migration")
	}
	// handle SERVER errors
	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server : %w", err)
	}

	return nil
}
