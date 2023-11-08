package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AbdelilahOu/GoThingy/config"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// connect to db
	testDb, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDb)

	os.Exit(m.Run())
}
