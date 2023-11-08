package main

import (
	"database/sql"
	"log"

	"github.com/AbdelilahOu/GoThingy/api"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// connect to db
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// create store and server
	store := db.NewStore(conn)
	server := api.NewServer(store)
	// start server
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
