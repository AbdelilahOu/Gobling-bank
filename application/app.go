package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	db     *sqlx.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mypassword"
	dbname   = "postgres"
)

func New() *App {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env doesn't exist or couldn't be loaded:", err)
	}

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
	//
	app := &App{
		db: db,
	}
	// Load routes
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	// Create a server
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	// Ping the database
	err := a.db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
	}
	// run migrations
	err = runDBMigrations(a.db)
	if err != nil {
		fmt.Println("Error running migrations:", err)
	}
	// close db on gracefull shudown
	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()
	// make channel to recieve events from go routine
	ch := make(chan error, 1)
	// gracefullshutdown
	go func() {
		err := server.ListenAndServe()
		// Handle server errors
		if err != nil {
			ch <- fmt.Errorf("error starting the server: %v", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

func runDBMigrations(db *sqlx.DB) error {
	// Run migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	// Specify the correct path to your migrations directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
