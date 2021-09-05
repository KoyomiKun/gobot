package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var (
	log_path string
)

func main() {
	app := &cli.App{
		Name:    "gobot",
		Usage:   "QQ Bot in golang",
		Version: "1.0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "config log file path",
				Value:   "./log",
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	app.Run(os.Args)
}
