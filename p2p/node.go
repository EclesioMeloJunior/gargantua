package p2p

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	badger "github.com/ipfs/go-ds-badger2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	initialTTLAdvertisementTimeout = time.Millisecond
	tryAdvertiseTimeout            = time.Second * 30
	findPeersTimeout               = time.Second * 10

	defaultDatastorePath = "libp2p-datastore"
)

type (
	Node struct {
		ds   *badger.Datastore
		pid  protocol.ID
		ctx  context.Context
		Host host.Host
		dht  *dual.DHT

		bootnodes []peer.AddrInfo
	}
)

// Close will close the libp2p host conns
func (n *Node) Close() error {
	return n.Host.Close()
}

// setupBootNodes will connect to bootnodes
func (n *Node) setupBootNodes() error {
	for _, peerinfo := range n.bootnodes {
		n.Host.Peerstore().AddAddrs(peerinfo.ID, peerinfo.Addrs, peerstore.PermanentAddrTTL)

		ttlctx, cancel := context.WithTimeout(n.ctx, time.Second*2)
		defer cancel()

		if err := n.Host.Connect(ttlctx, peerinfo); err != nil {
			log.Println("could not connect bootnode", peerinfo.ID)
			continue
		}

		log.Println("connected to bootnode", peerinfo.ID)
	}

	return nil
}

func (n *Node) discoveryAndAdvertise() error {
	if err := n.dht.Bootstrap(n.ctx); err != nil {
		return err
	}

	routingDiscovery := discovery.NewRoutingDiscovery(n.dht)
	discovery.Advertise(n.ctx, routingDiscovery, string(n.pid))

	ticker := time.NewTicker(findPeersTimeout)

	go func() {
		// log.Println("finding peers...")
		for {
			select {
			case <-n.ctx.Done():
				return
			case <-ticker.C:
				peers, err := discovery.FindPeers(n.ctx, routingDiscovery, string(n.pid))
				if err != nil {
					return
				}

				for _, peerinfo := range peers {
					if peerinfo.ID == n.Host.ID() || peerinfo.ID == "" {
						continue
					}

					// log.Println("found peer:", peerinfo.ID)
					if n.Host.Network().Connectedness(peerinfo.ID) != network.Connected {
						err := n.Host.Connect(n.ctx, peerinfo)
						if err != nil {
							log.Println("could not be able to connect to peer:", peerinfo.ID)
						} else {
							log.Println("connected to", peerinfo.Addrs)
							log.Println("connected peers", n.Host.Network().Peers())
						}
					}
				}
			}
		}
	}()

	log.Println("DHT discovery started")
	return nil
}

func (n *Node) MultiAddrs() (maddrs []ma.Multiaddr) {
	addrs := n.Host.Addrs()

	for _, addr := range addrs {
		maddr, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", addr, n.Host.ID()))
		if err != nil {
			continue
		}

		maddrs = append(maddrs, maddr)
	}

	return maddrs
}

func (n *Node) StartDiscovery() error {
	if len(n.bootnodes) < 1 {
		peers := n.Host.Network().Peers()

		for {
			if len(peers) > 0 {
				break
			}

			select {
			case <-time.After(time.Second * 5):
				log.Println("no peers yet, waiting connections...")
			case <-n.ctx.Done():
				return nil
			}

			peers = n.Host.Network().Peers()
		}

		for _, p := range peers {
			n.bootnodes = append(n.bootnodes, n.Host.Peerstore().PeerInfo(p))
		}
	}

	log.Println("starting DHT ...", n.bootnodes)
	dhtopts := []dual.Option{
		dual.DHTOption(kaddht.Datastore(n.ds)),
		dual.DHTOption(kaddht.BootstrapPeers(n.bootnodes...)),
		dual.DHTOption(kaddht.V1ProtocolOverride(n.pid + "/kad")),
		dual.DHTOption(kaddht.Mode(kaddht.ModeAutoServer)),
	}

	dht, err := dual.New(n.ctx, n.Host, dhtopts...)
	if err != nil {
		return err
	}

	n.dht = dht
	return n.discoveryAndAdvertise()
}

func NewP2PNode(ctx context.Context, pid protocol.ID, basepath, port string, bootnodes []string) (*Node, error) {
	bootnodesInfos, err := stringsToAddrInfo(bootnodes)
	if err != nil {
		log.Println("could not parse bootnode strings to addr")
	}

	opts, err := buildP2Popts(port, ctx)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	ds, err := badger.NewDatastore(path.Join(basepath, defaultDatastorePath), &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}

	n := &Node{
		ds:        ds,
		pid:       pid,
		ctx:       ctx,
		Host:      host,
		bootnodes: bootnodesInfos,
	}

	if err := n.setupBootNodes(); err != nil {
		return nil, err
	}

	return n, nil
}
