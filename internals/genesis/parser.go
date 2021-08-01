package genesis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	Genesis struct {
		ChainName   string            `json:"chainName"`
		Balances    map[string]uint32 `json:"balances"`
		Authorities []string          `json:"authorities"`
	}
)

func parseFromJSON(b []byte) (*Genesis, error) {
	g := new(Genesis)
	if err := json.Unmarshal(b, g); err != nil {
		return nil, err
	}

	return g, nil
}

func ReadGenesis(basepath, chain string) (*Genesis, error) {
	genesisfile := fmt.Sprintf("./chain/%s/genesis.json", chain)
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

	return parseFromJSON(genesisbytes)
}
