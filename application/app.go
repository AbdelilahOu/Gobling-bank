package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
	// defer db.Close()

	// Ping the database
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
	}

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

	// Run migrations
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		fmt.Println("Error getting driver:", err)
		return err
	}
	// migrations path
	migrationsPath, err := filepath.Rel("/", "./db/migrations")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return err
	}
	// Specify the correct path to your migrations directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		fmt.Println("Error creating migrations instance:", err)
		return err
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Error applying migrations:", err)
		return err
	}

	// Handle server errors
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
		return err
	}

	return nil
}
