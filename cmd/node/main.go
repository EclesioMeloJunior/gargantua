package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EclesioMeloJunior/gargantua/config"
	"github.com/EclesioMeloJunior/gargantua/p2p"
	"github.com/libp2p/go-libp2p-core/protocol"
)

const defaultConfigPath = "./config.dev.json"

func main() {
	configpath := flag.String("config", defaultConfigPath, "path to config json file")
	flag.Parse()

	nodeconfig, err := config.FromJson(*configpath)
	if err != nil {
		log.Println("problem to load config", err)
		return
	}

	expandedDir, err := config.ExpandDir(nodeconfig.Node.Basepath)
	if err != nil {
		log.Println(err)
		return
	}

	if err := config.SetupBasepath(expandedDir); err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n, err := p2p.NewP2PNode(ctx, protocol.ID(nodeconfig.Node.Protocol), expandedDir, nodeconfig.Network.Port, nodeconfig.Network.Bootnodes)
	if err != nil {
		log.Printf("could not start node: %v\n", err)
		return
	}

	log.Println("node started", n.Host.ID())
	log.Println("Addresses", n.MultiAddrs())

	if err := n.StartDiscovery(); err != nil {
		log.Println("failed start discovery", err)
		return
	}

	log.Println("protocols", n.Host.Mux().Protocols())

	rpcservice := p2p.NewRPC(n.Host, protocol.ID(nodeconfig.Node.Protocol))
	if err = rpcservice.Setup(); err != nil {
		log.Println("failed to setup rpc", err)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("shutting down...")
}
