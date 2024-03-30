package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AbdelilahOu/GoThingy/api"
	"github.com/AbdelilahOu/GoThingy/config"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/worker"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// connect to db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// start redis
	redisOps := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOps)
	store := db.NewStore(conn)
	go runTaskProcessor(redisOps, store)
	// create store and server
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		panic(err)
	}
	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	err := taskProcessor.Start()
	if err != nil {
		fmt.Println(err)
	}
}
