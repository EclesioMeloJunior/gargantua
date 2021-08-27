package block

import (
	"crypto/sha256"
	"fmt"
)

type Hash [32]byte

func (h *Hash) String() string {
	return fmt.Sprintf("0x%x", h[:])
}

type Hashable interface {
	Hash() (Hash, error)
}

func NewSHA256Hash(data []byte) Hash {
	return sha256.Sum256(data)
}
