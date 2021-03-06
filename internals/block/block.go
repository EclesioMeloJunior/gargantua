package block

import (
	"encoding/binary"
	"math/big"

	"github.com/EclesioMeloJunior/gargantua/internals/encoding"
)

type Body []byte

type Header struct {
	ParentHash Hash
	BlockHash  Hash
	StateRoot  Hash
	CreatedAt  *big.Int // timestamp
}

func NewHeader(parentHash Hash, root Hash, createdAt int64) *Header {
	var createdAtBytes [8]byte
	binary.LittleEndian.PutUint64(createdAtBytes[:], uint64(createdAt))

	toHash := make([]byte, 0)
	toHash = append(toHash, parentHash[:]...)
	toHash = append(toHash, root[:]...)
	toHash = append(toHash, createdAtBytes[:]...)

	return &Header{
		ParentHash: parentHash,
		StateRoot:  root,
		BlockHash:  NewSHA256Hash(toHash),
		CreatedAt:  big.NewInt(createdAt),
	}
}

// Hash generate the sha256 hash to then rlp encoded Header struct
func (h Header) Hash() (Hash, error) {
	// if blockhash is empty then update and return
	if len(h.BlockHash[:]) < 1 {
		generatedHash, err := encoding.RLPEncode(h)
		if err != nil {
			return Hash{}, err
		}

		h.BlockHash = NewSHA256Hash(generatedHash)
	}

	return h.BlockHash, nil
}

type Block struct {
	Body   Body
	Header *Header
}

func NewEmptyBlock() *Block {
	return &Block{
		Header: &Header{},
		Body:   []byte{},
	}
}
