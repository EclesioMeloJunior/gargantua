package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/EclesioMeloJunior/gargantua/p2p"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/urfave/cli/v2"
)

var NodeCmd = &cli.Command{
	Name:  "node",
	Usage: "setup a gargantua node",
	Subcommands: []*cli.Command{
		{
			Name:   "initialize",
			Usage:  "start a non-validator node by default",
			Action: initialize,
			Flags: append(globalFlags, &cli.StringFlag{
				Required: true,
				Name:     "key",
				Aliases:  []string{"k"},
			}),
		},
	},
}

func initialize(c *cli.Context) error {
	nodeconfig, err := config.FromJson(c.String("config"))
	if err != nil {
		return err
	}

	expandedDir, err := config.ExpandDir(nodeconfig.Node.Basepath)
	if err != nil {
		return err
	}

	hasKeyPair, err := checkNodeHasKeyPair(expandedDir, c.String("key"))
	if err != nil {
		return err
	}

	if !hasKeyPair {
		return errors.New("key pairs not found. execute gg key new --name={some-key-name} to create a new key pair")
	}

	if err := config.SetupBasepath(expandedDir); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n, err := p2p.NewP2PNode(ctx, protocol.ID(nodeconfig.Node.Protocol), expandedDir, nodeconfig.Network.Port, nodeconfig.Network.Bootnodes)
	if err != nil {
		return err
	}

	log.Println("node started", n.Host.ID())
	log.Println("Addresses", n.MultiAddrs())

	if err := n.StartDiscovery(); err != nil {
		return err
	}

	log.Println("protocols", n.Host.Mux().Protocols())

	rpcservice := p2p.NewRPC(n.Host, protocol.ID(nodeconfig.Node.Protocol))
	if err = rpcservice.Setup(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("shutting down...")

	return nil
}

func checkNodeHasKeyPair(basepath, name string) (bool, error) {
	keysdir, err := config.ExpandDir(filepath.Join(basepath, "keys"))
	if err != nil {
		return false, err
	}

	publicKeyPath := fmt.Sprintf(keystore.DefaultKeystoreFile, keysdir, name, keystore.PublicType)
	privateKeyPath := fmt.Sprintf(keystore.DefaultKeystoreFile, keysdir, name, keystore.PrivateType)

	pubExists, err := checkFileStat(publicKeyPath)
	if err != nil {
		return false, err
	}

	privExists, err := checkFileStat(privateKeyPath)
	if err != nil {
		return false, err
	}

	return pubExists && privExists, nil
}

func checkFileStat(filepath string) (bool, error) {
	finfo, err := os.Stat(filepath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, errors.New("node doesnt have key pair")
	}

	return finfo != nil, nil
}
