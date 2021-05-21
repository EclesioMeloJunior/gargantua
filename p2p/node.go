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

func (n *Node) setupBootstrapPeers() {
	for _, addr := range kaddht.DefaultBootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Println("error to get the addr infor from peer", addr)
			continue
		}

		if err := n.Host.Connect(n.ctx, *peerinfo); err != nil {
			log.Println("could not connect to peer", peerinfo.Addrs)
			continue
		}

		n.bootnodes = append(n.bootnodes, n.Host.Peerstore().PeerInfo(peerinfo.ID))
		log.Println("connected to bootstrap peer", addr)
	}
}

func (n *Node) discoveryAndAdvertise() error {
	routingDiscovery := discovery.NewRoutingDiscovery(n.dht)
	err := n.dht.Bootstrap(n.ctx)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	log.Println("current peers:", len(n.Host.Network().Peers()))

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

func NewP2PNode(ctx context.Context, pid protocol.ID, basepath, port string, bootnodes []string) (*Node, error) {
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
		ctx:  ctx,
		Host: host,
	}

	n.setupBootstrapPeers()
	dhtopts := []dual.Option{
		dual.DHTOption(kaddht.Datastore(ds)),
		dual.DHTOption(kaddht.V1ProtocolOverride(pid + "/kad")),
		dual.DHTOption(kaddht.Mode(kaddht.ModeAutoServer)),
		dual.DHTOption(kaddht.BootstrapPeers(n.bootnodes...)),
	}

	dht, err := dual.New(ctx, host, dhtopts...)
	if err != nil {
		return nil, err
	}

	n.dht = dht
	if err := n.discoveryAndAdvertise(); err != nil {
		return nil, err
	}

	return n, nil
}
