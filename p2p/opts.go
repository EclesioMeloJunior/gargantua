package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	secio "github.com/libp2p/go-libp2p-secio"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

func buildP2Popts(port string, ctx context.Context) ([]libp2p.Option, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.ListenAddrStrings(getListenAddrs(port)...),
		libp2p.ConnectionManager(connmgr.NewConnManager(100, 400, time.Minute)),
		libp2p.NATPortMap(),
		libp2p.DisableRelay(),
	}

	return opts, nil
}
