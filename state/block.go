package state

type (
	BlockHeader struct {
		Previus Hash
	}

	Block struct {
		Header BlockHeader
	}
)

func (b *Block) Hash() (Hash, error) {
	return Hash{}, nil
}
