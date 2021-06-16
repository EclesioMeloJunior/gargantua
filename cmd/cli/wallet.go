package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/EclesioMeloJunior/gargantua/wallet"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var WalletCmd = &cli.Command{
	Name:  "wallet",
	Usage: "manage wallets private and public keys",
	Subcommands: []*cli.Command{
		{
			Name:   "new",
			Usage:  "create and stores the keypair",
			Action: newWalletAddress,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Required: true,
					Name:     "name",
					Aliases:  []string{"n"},
				},
			},
		},

		{
			Name:   "addresses",
			Usage:  "list all address in the current node",
			Action: listAddresses,
		},
	},
}

func newWalletAddress(c *cli.Context) error {
	nodeconfig, err := config.FromJson(c.String("config"))
	basepath, err := config.ExpandDir(nodeconfig.Node.Basepath)
	if err != nil {
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

	walletdir, err := config.ExpandDir(filepath.Join(basepath, "wallets"))
	if err != nil {
		return err
	}

	fmt.Println(walletdir)

	err = keystore.StoreKeyPair(c.String("name"), walletdir, pair, password)
	if err != nil {
		return err
	}

	addr := wallet.GetAddress(pair.Public)
	fmt.Printf("\nAddress: %s\n", addr)
	return nil
}

func listAddresses(c *cli.Context) error {
	return errors.New("not implemented yet")
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

	cmp := strings.Compare(
		strings.TrimSpace(string(bytePassword)),
		strings.TrimSpace(string(byteConfirmPassword)),
	)

	if cmp != 0 {
		return "", fmt.Errorf("\npasswords not match")
	}

	return string(bytePassword), nil
}
