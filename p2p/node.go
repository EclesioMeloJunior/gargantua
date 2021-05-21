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
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
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

func (n *Node) setupBootnodes() {
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
}

func (n *Node) waitPeersToDHT() {
	if len(n.bootnodes) == 0 {
		peers := n.Host.Network().Peers()

		for {
			if len(peers) > 0 {
				break
			}

			select {
			case <-time.After(time.Second * 10):
				log.Println("no peers yet, waiting peers to start dht")
			case <-n.ctx.Done():
				return
			}

			peers = n.Host.Network().Peers()
		}

		for _, p := range peers {
			n.bootnodes = append(n.bootnodes, n.Host.Peerstore().PeerInfo(p))
		}
	}
}

func (n *Node) discoveryAndAdvertise() error {
	routingDiscovery := discovery.NewRoutingDiscovery(n.dht)

	err := n.dht.Bootstrap(n.ctx)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	go func() {
		ttl := initialTTLAdvertisementTimeout
		for {
			select {
			case <-time.After(ttl):
				log.Println("advertising ourselves in the DHT...")
				err := n.dht.Bootstrap(n.ctx)
				if err != nil {
					log.Println("failed to bootstrap DHT")
					continue
				}

				ttl, err = routingDiscovery.Advertise(n.ctx, string(n.pid))
				if err != nil {
					log.Println("fail to advertise in the DHT", err)
					ttl = tryAdvertiseTimeout
				}
			case <-n.ctx.Done():
				return
			}
		}
	}()

	go func() {
		log.Println("finding peers...")
		peerC, err := routingDiscovery.FindPeers(n.ctx, string(n.pid))
		if err != nil {
			log.Println("could not initialize to find peers", err)
			return
		}

		for {
			select {
			case <-n.ctx.Done():
				return
			case <-time.After(findPeersTimeout):
				log.Println("check current peers amount...")
			case peerinfo := <-peerC:
				if peerinfo.ID == n.Host.ID() || peerinfo.ID == "" {
					continue
				}
				log.Println("found peer:", peerinfo.ID)

				err := n.Host.Connect(n.ctx, peerinfo)
				if err != nil {
					log.Println("could not be able to connect to peer:", peerinfo.ID)
					continue
				}

				log.Println("connected to", peerinfo.Addrs)
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
	n.setupBootnodes()
	n.waitPeersToDHT()

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
	var err error

	if err != nil {
		return nil, err
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

	bootnodesInfos, err := stringsToAddrInfo(bootnodes)
	if err != nil {
		log.Println("could not parse bootnode strings to addr")
	}

	n := &Node{
		ds:        ds,
		pid:       pid,
		ctx:       ctx,
		Host:      host,
		bootnodes: bootnodesInfos,
	}

	return n, nil
}
