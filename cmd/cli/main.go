package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

const defaultConfigPath = "./config.dev.json"

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		HasBeenSet: false,
		Name:       "config",
		Aliases:    []string{"c"},
		Value:      defaultConfigPath,
	},
}

func main() {
	app := &cli.App{
		Name:  "gg",
		Usage: "Gargantua",
	}

	app.Commands = []*cli.Command{
		NodeCmd,
		KeysCmd,
	}

	// global flags
	app.Flags = globalFlags

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
