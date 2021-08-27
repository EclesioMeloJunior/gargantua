package p2p

import (
	"context"
	"log"
	"reflect"

	"github.com/EclesioMeloJunior/gargantua/p2p/noderpc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
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
	blockHdl := new(noderpc.BlockHandler)

	rpc.server = libp2prpc.NewServer(rpc.node, rpc.pid)
	rpc.server.Register(blockHdl)

	// registry p2p rpc handlers
	rpc.client = libp2prpc.NewClient(rpc.node, rpc.pid)
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

func (rpc *RPC) Call(service, function string, args, reply interface{}) ([]interface{}, []error) {
	errs := make([]error, 0)
	peers := rpc.node.Network().Peers()

	replies := make([]interface{}, 0)
	for range peers {
		ptr := reflect.New(reflect.TypeOf(reply))
		replies = append(replies, ptr)
	}

	for i, p := range peers {
		log.Println("sending to peer: ", peer.Encode(p))
		err := rpc.client.Call(p, service, function, args, replies[i])
		if err != nil {
			errs = append(errs, err)
		}
	}

	return replies, errs
}

func (rpc *RPC) Peers() peer.IDSlice {
	return rpc.node.Peerstore().Peers()
}

func contexts(n int) []context.Context {
	cxs := make([]context.Context, n)
	for i := 0; i < n; i++ {
		cxs[i] = context.Background()
	}
	return cxs
}
