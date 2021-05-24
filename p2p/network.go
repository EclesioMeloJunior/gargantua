package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	libp2prpc "github.com/libp2p/go-libp2p-gorpc"
)

type (
	RPC struct {
		server *libp2prpc.Server
		client *libp2prpc.Client
		pid    protocol.ID
		node   host.Host
	}
)

func NewRPC(node host.Host, pid protocol.ID) *RPC {
	return &RPC{
		pid:  pid,
		node: node,
	}
}

func (rpc *RPC) Setup() error {
	rpc.server = libp2prpc.NewServer(rpc.node, rpc.pid)
	// registry p2p rpc handlers
	rpc.client = libp2prpc.NewClientWithServer(rpc.node, rpc.pid, rpc.server)
	return nil
}

func (rpc *RPC) MultiCall(service, function string, args interface{}, replies []interface{}) []error {
	peers := rpc.node.Peerstore().Peers()

	return rpc.client.MultiCall(
		contexts(len(peers)),
		peers,
		service,
		function,
		args,
		replies,
	)
}

func contexts(n int) []context.Context {
	cxs := make([]context.Context, n)
	for i := 0; i < n; i++ {
		cxs[i] = context.Background()
	}
	return cxs
}
