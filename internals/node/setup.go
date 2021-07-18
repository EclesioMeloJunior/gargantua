package node

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	GenesisGeneralBucket  = []byte("general")
	GenesisBalancesBucket = []byte("balances")
	GenesisChainNameKey   = []byte("genesis_chain_name")
)

type Storage interface {
	HasKey(key []byte) (bool, error)
	StoreOnBucket(bucket, key, value []byte) error
}

func ReadGenesis(basepath, chain string) (*Genesis, error) {
	genesisfile := filepath.Join(basepath, fmt.Sprintf("%s.json", chain))
	_, err := os.Stat(genesisfile)
	if errors.Is(os.ErrNotExist, err) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	genesisbytes, err := ioutil.ReadFile(genesisfile)
	if err != nil {
		return nil, err
	}

	var genesis *Genesis
	if err = json.Unmarshal(genesisbytes, genesis); err != nil {
		return nil, err
	}

	return genesis, nil
}

func LoadGenesisOnStorage(g *Genesis, s Storage) error {
	// store chain name
	if err := tryStore(GenesisGeneralBucket, GenesisChainNameKey, []byte(g.Chain), s); err != nil {
		return err
	}

	// store balances
	for addr, value := range g.Balances {
		var bvalue [4]byte
		binary.LittleEndian.PutUint32(bvalue[:], value)

		addrBytes, err := hex.DecodeString(addr)
		if err != nil {
			return err
		}

		if err := tryStore(GenesisBalancesBucket, addrBytes, bvalue[:], s); err != nil {
			return err
		}
	}

	return nil
}

func tryStore(bucket, key, v []byte, s Storage) error {
	ok, err := s.HasKey(key)
	if err != nil {
		return err
	}

	if !ok {
		return s.StoreOnBucket(bucket, key, v)
	}

	return nil
}
