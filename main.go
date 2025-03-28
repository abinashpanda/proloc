package main

import (
	"fmt"
	"os"

	"github.com/abinashpanda/proloc/proloc"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "proloc",
		Description: "a command line utitlity to calculate the number of lines of code",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Usage:    "path of the project",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "ignore",
				Usage: "glob patterns to ignore",
			},
			&cli.Uint64Flag{
				Name:  "max-depth",
				Usage: "print the total for the directories (or file) only if is N or fewer levels",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := proloc.ProlocConfig{
				Project:  ctx.String("project"),
				Ignore:   ctx.StringSlice("ignore"),
				MaxDepth: ctx.Uint64("max-depth"),
			}
			return proloc.CountLines(config)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
