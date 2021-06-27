package proposal

import "github.com/EclesioMeloJunior/gargantua/internals/document"

type Proposal struct {
	Document   *document.Document `rlp:"document"`
	Signatures [][]byte           `rlp:"authors"`
}

func NewEmptyProposal() *Proposal {
	return &Proposal{
		Document:   nil,
		Signatures: [][]byte{},
	}
}
