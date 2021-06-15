package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

const defaultConfigPath = "./config.dev.json"

func main() {
	app := &cli.App{
		Name:  "gg",
		Usage: "Gargantua",
	}

	app.Commands = []*cli.Command{
		NodeCmd,
		WalletCmd,
	}

	// global flags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   defaultConfigPath,
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
