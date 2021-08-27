package noderpc

import (
	"context"
	"log"

	"github.com/EclesioMeloJunior/gargantua/internals/block"
	"github.com/EclesioMeloJunior/gargantua/internals/encoding"
)

// BlockHandler is the RPC handler for block related calls
type BlockHandler struct{}

func (h *BlockHandler) NewBlock(ctx context.Context, args []byte, res []byte) error {
	log.Println("block received")

	var b *block.Block
	err := encoding.RLPDecode(args, b)
	if err != nil {
		log.Printf("problem while decoding received block: %v\n", err)
		return err
	}

	bhash, err := b.Header.Hash()
	if err != nil {
		log.Printf("problem while hashing received block: %v\n", err)
		return err
	}

	log.Printf("block decoded: %s\n", bhash)
	return nil
}
