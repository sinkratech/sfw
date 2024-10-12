package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/samuelsih/formattag/align"
	"github.com/urfave/cli/v2"
)

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
