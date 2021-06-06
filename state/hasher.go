package state

import "encoding/hex"

type (
	Hash [32]byte

	Hasher interface {
		Hash() (Hash, error)
	}
)

func (h *Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h *Hash) String() string {
	return h.Hex()
}
