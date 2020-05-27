package main

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = []*cli.Command{
		{
			Name:  "new",
			Usage: "Creates a new project",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return errors.New("invalid argument")
				}

				config.Type = TypeWeekly // default
				if c.Bool("month") {
					config.Type = TypeMonthly
				}
				if c.Bool("day") {
					config.Type = TypeDaily
				}
				return config.New(c.Args().First())
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "month", Aliases: []string{"m"}},
				&cli.BoolFlag{Name: "week", Aliases: []string{"w"}},
				&cli.BoolFlag{Name: "day", Aliases: []string{"d"}},
			},
		},
		{
			Name:  "post",
			Usage: "Creates a new post",
			Action: func(c *cli.Context) error {
				if err := config.Load("config.yaml"); err != nil {
					return err
				}
				return config.Post()
			},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "template", Aliases: []string{"t"}},
			},
		},
		{
			Name:  "serve",
			Usage: "Builds html in temporary directory and runs local server",
			Action: func(c *cli.Context) error {
				if err := config.Load("config.yaml"); err != nil {
					return err
				}
				return config.Serve()
			},
		},
		{
			Name:  "pub",
			Usage: "Builds html in public directory",
			Action: func(c *cli.Context) error {
				if err := config.Load("config.yaml"); err != nil {
					return err
				}
				return config.Publish()
			},
		},
	}
}
