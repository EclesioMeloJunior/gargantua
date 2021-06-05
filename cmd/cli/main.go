package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gg",
		Usage: "Gargantua",
	}

	app.Commands = []*cli.Command{
		NodeCmd,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
