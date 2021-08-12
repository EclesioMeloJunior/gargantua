package block

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
		EthereumDB() ethdb.KeyValueStore
	}
)

func NewBlockFromGenesis(g *genesis.Genesis, s Storage) (*Block, error) {
	b := NewEmptyBlock()

	trieHash, err := trieFromGenesis(g, s.EthereumDB())
	if err != nil {
		return nil, err
	}

	b.Header = NewHeader(Hash{}, trieHash, time.Now().Unix())
	s.PutOnBucket(BlocksHashBucket, b.Header.BlockHash[:], trieHash[:])
	return b, nil
}

func trieFromGenesis(g *genesis.Genesis, db ethdb.KeyValueStore) (Hash, error) {
	trieDB := trie.NewDatabase(db)
	t, err := trie.New(common.Hash{}, trieDB)
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

		t.Update(accBytes, bvalue[:])
	}

	for _, auth := range g.Authorities {
		filledPattern := fmt.Sprintf(authoritiesTriePattern, auth)

		var bvalue [4]byte
		binary.LittleEndian.PutUint32(bvalue[:], 0)

		t.Update([]byte(filledPattern), bvalue[:])
	}

	t.Update(genesis.GenesisChainNameKey, []byte(g.ChainName))

	hash, err := t.Commit(nil)
	if err != nil {
		return Hash{}, err
	}

	if err = trieDB.Commit(hash, false, nil); err != nil {
		return Hash{}, nil
	}

	return Hash(hash), nil
}
