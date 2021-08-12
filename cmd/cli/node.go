package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
	"github.com/EclesioMeloJunior/gargantua/internals/block"
	"github.com/EclesioMeloJunior/gargantua/internals/genesis"
	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/EclesioMeloJunior/gargantua/p2p"
	"github.com/EclesioMeloJunior/gargantua/storage"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/urfave/cli/v2"
)

var NodeCmd = &cli.Command{
	Name:  "node",
	Usage: "setup a gargantua node",
	Subcommands: []*cli.Command{
		{
			Name:   "init",
			Usage:  "start a non-validator node by default",
			Action: initialize,
			Flags: append(globalFlags, &cli.StringFlag{
				Required: true,
				Name:     "key",
				Aliases:  []string{"k"},
			}, &cli.StringFlag{
				Required: true,
				Name:     "chain",
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

	hasKeyPair, err := keystore.CheckNodeHasKeyPair(expandedDir, c.String("key"))
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

	go n.StartDiscovery()

	log.Println("protocols", n.Host.Mux().Protocols())

	rpcservice := p2p.NewRPC(n.Host, protocol.ID(nodeconfig.Node.Protocol))
	if err = rpcservice.Setup(); err != nil {
		return err
	}

	storage, err := storage.NewStorage(expandedDir)
	if err != nil {
		return err
	}

	gn, err := genesis.ReadGenesis(expandedDir, c.String("chain"))
	if err != nil {
		return err
	}

	b, err := block.NewBlockFromGenesis(gn, storage)
	if err != nil {
		return err
	}

	log.Printf("genesis created: 0x%x", b.Header.BlockHash[:])

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("shutting down...")

	return nil
}
