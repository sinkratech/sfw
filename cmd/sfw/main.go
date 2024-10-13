package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/samuelsih/formattag/align"
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

func handleFormatTag(c *cli.Context) error {
	rootTarget := "."
	basefs := os.DirFS(".")

	if c.Args().First() != "" {
		rootTarget = c.Args().First()
	}

	err := fs.WalkDir(basefs, rootTarget, func(p string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		align.Init(p)

		b, errWalk := align.Do()
		if errWalk != nil {
			return fmt.Errorf("format align failed for %s: %w", p, errWalk)
		}

		errWalk = os.WriteFile(p, b, 0)
		if errWalk != nil {
			return fmt.Errorf("cannot write to file %s: %w", p, errWalk)
		}

		return nil
	})

	return err
}
