package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const name = "gournal"

var app = cli.NewApp()

func main() {
	os.Exit(run(os.Args))
}

const (
	exitOK = iota
	exitNG
)

func run(args []string) int {
	app.Name = name
	app.Usage = "Journal tool written in Go"
	app.Version = "0.0.1"

	if err := app.Run(args); err != nil {
		fmt.Println(err)
		return exitNG
	}
	return exitOK
}
