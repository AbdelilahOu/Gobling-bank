package main

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/GoThingy/api"
	"github.com/AbdelilahOu/GoThingy/config"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/mail"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/AbdelilahOu/GoThingy/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := utils.NewLogger()
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal("cannot load config:", err)
	}
	// connect to db
	pgPoom, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		logger.Fatal("cannot connect to db:", err)
	}
	// start redis
	redisOps := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOps)
	store := db.NewStore(pgPoom)
	go runTaskProcessor(config, redisOps, store)
	// create store and server
	server, err := api.NewServer(config, store, taskDistributor, *logger)
	if err != nil {
		panic(err)
	}
	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		logger.Fatal("cannot start server:", err)
	}
}

func runTaskProcessor(config config.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	err := taskProcessor.Start()
	if err != nil {
		fmt.Println(err)
	}
}
