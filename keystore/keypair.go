package keystore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type (
	Pair struct {
		Private *PrivateKey
		Public  *PublicKey
	}

	PrivateKey struct {
		*ecdsa.PrivateKey
	}

	PublicKey struct {
		*ecdsa.PublicKey
	}
)

func NewKeyPair() (*Pair, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Pair{
		Private: &PrivateKey{privKey},
		Public:  &PublicKey{&privKey.PublicKey},
	}, nil
}
