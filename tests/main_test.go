package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/AbdelilahOu/GoThingy/config"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

var testStore db.Store

func TestMain(m *testing.M) {
	var err error
	config, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// connect to db
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testStore = db.NewStore(connPool)
	os.Exit(m.Run())
}
