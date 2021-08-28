package p2p

import (
	"fmt"
	"log"

	libp2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

func (n *Node) SetupStreamHandlers(basepid string) {
	n.Host.SetStreamHandler(protocol.ID(fmt.Sprintf("%s/arx", basepid)), func(s libp2pnetwork.Stream) {
		log.Println("/arx received!")
		s.Close()
	})
}
