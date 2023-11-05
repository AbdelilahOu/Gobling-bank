package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDb)

	os.Exit(m.Run())
}
