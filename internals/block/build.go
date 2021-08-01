package block

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/EclesioMeloJunior/gargantua/internals/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

const (
	authoritiesTriePattern = "authoritie:%s"
)

var (
	BlocksHashBucket = []byte("blocks_hash")
)

type (
	Storage interface {
		PutOnBucket(bucket, key, value []byte) error
		Put(key, value []byte) error
		Delete(key []byte) error

		ethdb.KeyValueStore

		io.Closer
	}
)

func NewBlockFromGenesis(g *genesis.Genesis, s Storage) (*Block, error) {
	b := NewEmptyBlock()

	trieHash, err := trieFromGenesis(g, s)
	if err != nil {
		return nil, err
	}

	b.Header = NewHeader(Hash{}, trieHash, time.Now().Unix())
	s.PutOnBucket(BlocksHashBucket, b.Header.BlockHash[:], trieHash[:])
	fmt.Printf("root state from block: %x\n", trieHash[:])
	return b, nil
}

func trieFromGenesis(g *genesis.Genesis, s Storage) (Hash, error) {
	trie, err := trie.New(common.Hash{}, trie.NewDatabase(s))
	if err != nil {
		return Hash{}, err
	}

	for acc, bal := range g.Balances {
		accBytes, err := hex.DecodeString(acc[2:])
		if err != nil {
			return Hash{}, err
		}

		var bvalue [4]byte
		binary.LittleEndian.PutUint32(bvalue[:], bal)

		trie.Update(accBytes, bvalue[:])
	}

	for _, auth := range g.Authorities {
		filledPattern := fmt.Sprintf(authoritiesTriePattern, auth)

		var bvalue [4]byte
		binary.LittleEndian.PutUint32(bvalue[:], 0)

		trie.Update([]byte(filledPattern), bvalue[:])
	}

	trie.Update(genesis.GenesisChainNameKey, []byte(g.ChainName))
	hash, err := trie.Commit(nil)

	if err != nil {
		return Hash{}, err
	}

	return Hash(hash), nil
}
