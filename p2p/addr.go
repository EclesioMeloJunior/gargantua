package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func getListenAddrs(port string) []string {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/" + port,
	}

	return addrs
}

func stringToAddrInfo(s string) (*peer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(s)
	if err != nil {
		return nil, err
	}

	return peer.AddrInfoFromP2pAddr(maddr)
}

func stringsToAddrInfo(s []string) ([]peer.AddrInfo, error) {
	pinfos := make([]peer.AddrInfo, len(s))
	for i, v := range s {
		p, err := stringToAddrInfo(v)
		if err != nil {
			return nil, err
		}

		pinfos[i] = *p
	}

	return pinfos, nil
}
