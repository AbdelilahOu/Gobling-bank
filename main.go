package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/AbdelilahOu/GoThingy/application"
)

func main() {
	// create app instance
	app := application.New()
	// use signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	// start app
	err := app.Start(ctx)
	// app errors
	if err != nil {
		fmt.Println("failed to start app : %w", err)
	}
}
