package keystore

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func GetAddress(pk *PublicKey) string {
	pkBytes := crypto.FromECDSAPub(pk.PublicKey)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pkBytes[1:])
	return hexutil.Encode(hash.Sum(nil)[12:])
}
