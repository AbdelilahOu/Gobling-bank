package main

import (
	"fmt"
	"context"
	"github.com/AbdelilahOu/GoThingy/application"
)

func main(){
	app := application.New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start app : %w",err)
	}
}