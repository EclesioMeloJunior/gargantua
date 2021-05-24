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

const defaultBasepath = "~/.gargantua"
const defaultPID = protocol.ID("/gargantua/dev/v0")

var (
	port     string
	basepath string
	bootnode string
)

func init() {
	flag.StringVar(&port, "port", "9002", "setup the port to listen to")
	flag.StringVar(&basepath, "basepath", defaultBasepath, "the directory to stores node related files")
	flag.StringVar(&bootnode, "bootnode", "", "setup a bootnode as a peer")
}

func main() {
	flag.Parse()

	expandedDir, err := config.ExpandDir(basepath)
	if err != nil {
		log.Println(err)
		return
	}

	if err := config.SetupBasepath(expandedDir); err != nil {
		log.Println(err)
		return
	}

	bootnodes := []string{}
	if bootnode != "" {
		bootnodes = append(bootnodes, bootnode)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n, err := p2p.NewP2PNode(ctx, defaultPID, expandedDir, port, bootnodes)
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

	rpcservice := p2p.NewRPC(n.Host, defaultPID)
	if err = rpcservice.Setup(); err != nil {
		log.Println("failed to setup rpc", err)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("shutting down...")
}
