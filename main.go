package main

import (
	"log"
	"os"

	"github.com/Meland-Inc/logservice/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "serve",
		Usage: "Start dapr miro service",
		Action: func(c *cli.Context) error {
			return cmd.Serve(c)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
