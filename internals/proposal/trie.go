package proposal

type ProposalTrie interface {
	Insert()
	Lookup(p []byte) *Proposal
}

func CreateProposalTrie([]*Proposal) {

}
