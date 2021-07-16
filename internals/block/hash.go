package block

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/rlp"
)

type Hash [32]byte

type Hashable interface {
	Hash() (Hash, error)
}

func NewSHA256Hash(data []byte) Hash {
	return sha256.Sum256(data)
}

func RLPAndSHA256Hash(v interface{}) (Hash, error) {
	hb, err := rlp.EncodeToBytes(v)

	if err != nil {
		return Hash{}, err
	}

	return NewSHA256Hash(hb), nil
}
