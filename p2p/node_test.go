package p2p_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/EclesioMeloJunior/gargantua/p2p"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/assert"
)

func TestNodeDiscovery(t *testing.T) {
	ctx := context.Background()

	nodesCount := 3
	basepaths := make([]string, nodesCount)
	for i := 0; i < nodesCount; i++ {
		basepath, err := ioutil.TempDir("", "*")
		assert.NoError(t, err)
		basepaths[i] = basepath
	}

	defer func() {
		for _, path := range basepaths {
			assert.NoError(t, os.RemoveAll(path))
		}
	}()

	// creates node A using the basepath 0
	nodeA, err := p2p.NewP2PNode(ctx, protocol.ID("testing"), basepaths[0], "9001", []string{})
	assert.NoError(t, err)
	assert.NotEmpty(t, nodeA.Host.Addrs())

	// creates node B using the basepath 1
	// and using node A as bootnode
	nodeB, err := p2p.NewP2PNode(ctx, protocol.ID("testing"), basepaths[1], "9002", []string{
		nodeA.MultiAddrs()[0].String(),
	})
	assert.NoError(t, err)

	// creates node C using the basepath 2
	// and using node A as bootnode
	nodeC, err := p2p.NewP2PNode(ctx, protocol.ID("testing"), basepaths[2], "9003", []string{
		nodeA.MultiAddrs()[0].String(),
	})
	assert.NoError(t, err)

	nodeA.StartDiscovery()
	nodeB.StartDiscovery()
	nodeC.StartDiscovery()

	testCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ticker := time.NewTicker(time.Second * 2)

TestLoop:
	for {
		select {
		case <-ticker.C:
			nodeCAddrInfo := nodeB.Host.Peerstore().PeerInfo(nodeC.Host.ID())
			if assert.NotNil(t, nodeCAddrInfo) {
				cancel()
			}
		case <-testCtx.Done():
			ticker.Stop()
			break TestLoop
		}
	}

}
