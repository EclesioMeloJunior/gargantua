package genesis

import (
	"encoding/binary"
	"encoding/hex"
)

var (
	GenesisGeneralBucket  = []byte("general")
	GenesisBalancesBucket = []byte("balances")
)

var (
	GenesisChainNameKey = []byte("genesis_chain_name")
)

type Storage interface {
	PutOnBucket(bucket, key, value []byte) error
}

func tryStoreGenesisInfo(bucket, key, v []byte, s Storage) error {
	return s.PutOnBucket(bucket, key, v)
}

func StoreGenesis(g *Genesis, s Storage) error {
	// store chain name
	if err := tryStoreGenesisInfo(GenesisGeneralBucket, GenesisChainNameKey, []byte(g.ChainName), s); err != nil {
		return err
	}

	// store balances
	for acc, value := range g.Balances {
		var bvalue [4]byte
		binary.LittleEndian.PutUint32(bvalue[:], value)

		addrBytes, err := hex.DecodeString(acc[2:])
		if err != nil {
			return err
		}

		if err := tryStoreGenesisInfo(GenesisBalancesBucket, addrBytes, bvalue[:], s); err != nil {
			return err
		}
	}
	return nil
}
