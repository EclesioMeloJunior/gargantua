package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	GenesisStorageChain = "genesis_chain_name"
)

type Storage interface {
	HasKey(key []byte) (bool, error)
	Store(key []byte, value []byte) error
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
	if ok, err := s.HasKey([]byte(g.Chain)); err != nil {
		return err
	} else if !ok {
		return s.Store([]byte(GenesisStorageChain), []byte(g.Chain))
	}

	return nil
}
