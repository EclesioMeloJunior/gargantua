package encoding

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func RLPDecode(b []byte, v interface{}) error {
	return rlp.DecodeBytes(b, v)
}

func RLPEncode(v interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(v)
}
