package main

import (
	"os"

	"github.com/sinkratech/sfw/cmd/sfw/feature"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Suggest: true,
		Name:    "sfw",
		Usage:   "Mini framework for development with huma",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"gen", "g"},
				Usage:   "Generate something you want (features, i18n, etc.)",
				Subcommands: []*cli.Command{
					{
						Name:    "feature",
						Aliases: []string{"feat", "f"},
						Usage:   "Generate features in target directory",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "base",
								Value: "api",
								Usage: "Base directory for all features",
							},
						},
						Action: feature.Handle,
					},
				},
			},
			{
				Name:    "formattag",
				Aliases: []string{"ft", "formatag"},
				Usage:   "Format struct tag",
				Action:  handleFormatTag,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}
