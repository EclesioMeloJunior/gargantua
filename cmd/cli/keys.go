package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var KeysCmd = &cli.Command{
	Name:  "key",
	Usage: "manage private and public keys",
	Subcommands: []*cli.Command{
		{
			Name:   "new",
			Usage:  "create and stores the keypair",
			Action: newKeyPairCmd,
			Flags: append(globalFlags, &cli.StringFlag{
				Required: true,
				Name:     "name",
				Aliases:  []string{"n"},
			}),
		},

		{
			Name:   "address",
			Usage:  "get the public address of a specific key",
			Action: getPublicAddress,
			Flags: append(globalFlags, &cli.StringFlag{
				Required: true,
				Name:     "key",
				Aliases:  []string{"k"},
			}),
		},
	},
}

func newKeyPairCmd(c *cli.Context) error {
	nodeconfig, err := config.FromJson(c.String("config"))
	if err != nil {
		return err
	}

	basepath, err := config.ExpandDir(nodeconfig.Node.Basepath)
	if err != nil {
		return err
	}

	if err = config.SetupBasepath(basepath); err != nil {
		return err
	}

	password, err := readPassword()
	if err != nil {
		return err
	}

	pair, err := keystore.NewKeyPair()
	if err != nil {
		return err
	}

	keysdir, err := config.ExpandDir(filepath.Join(basepath, "keys"))
	if err != nil {
		return err
	}

	err = keystore.StoreKeyPair(c.String("name"), keysdir, pair, password)
	if err != nil {
		return err
	}

	addr := keystore.GetAddress(pair.Public)
	fmt.Printf("\nAddress: %s\n", addr)
	return nil
}

func getPublicAddress(c *cli.Context) error {
	nodeconfig, err := config.FromJson(c.String("config"))
	if err != nil {
		return err
	}

	var basepath string
	if basepath, err = config.ExpandDir(nodeconfig.Node.Basepath); err != nil {
		return err
	}

	pubkey, err := keystore.LoadPublicKey(basepath, c.String("key"))
	if err != nil {
		return err
	}

	addr := keystore.GetAddress(pubkey)
	fmt.Printf("Address: %s\n", addr)
	return nil
}

func readPassword() (string, error) {
	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	if len(bytePassword) == 0 {
		return "", fmt.Errorf("\nempty password")
	}

	fmt.Printf("\nConfirm password: ")

	byteConfirmPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Print("\n")

	cmp := strings.Compare(
		strings.TrimSpace(string(bytePassword)),
		strings.TrimSpace(string(byteConfirmPassword)),
	)

	if cmp != 0 {
		return "", fmt.Errorf("\npasswords not match")
	}

	return string(bytePassword), nil
}
