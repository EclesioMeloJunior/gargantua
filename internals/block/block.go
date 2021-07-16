package block

import (
	"bytes"
)

type Body struct {
	Transactions []Transaction
}

type Header struct {
	ParentHash Hash
	BlockHash  Hash
	TxRoot     Hash
	CreatedAt  int64 // timestamp
}

// Hash generate the sha256 hash to then rlp encoded Header struct
func (h Header) Hash() (Hash, error) {
	// if blockhash is empty then update and return
	if bytes.Equal(h.BlockHash[:], []byte{}) {
		generatedHash, err := RLPAndSHA256Hash(h)
		if err != nil {
			return Hash{}, err
		}

		h.BlockHash = generatedHash
	}

	return h.BlockHash, nil
}

type Block struct {
	Body []byte
}
