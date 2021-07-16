package keystore

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type Address [32]byte

func (a *Address) String() string {
	return hexutil.Encode(a[12:])
}

func GetAddress(pk *PublicKey) Address {
	pkBytes := crypto.FromECDSAPub(pk.PublicKey)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pkBytes[1:])

	var addr Address
	copy(addr[:], hash.Sum(nil)[:])

	return addr
}
