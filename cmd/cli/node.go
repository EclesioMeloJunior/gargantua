package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
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
