package main

import (
	"log"

	"textstat/pkg/cli"
)

func main() {
	app := cli.NewApp()

	if err := app.Run(); err != nil {
		app.HandleError(err)
		log.Fatal("Application failed to run")
	}
}
