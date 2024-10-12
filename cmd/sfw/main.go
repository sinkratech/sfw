package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sfw",
		Usage: "Mini framework for development with huma",
		Commands: []*cli.Command{
			{
				Name:    "formattag",
				Aliases: []string{"ft", "formatag"},
				Usage:   "Format struct tag",
				Action:  handleFormatTag,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
