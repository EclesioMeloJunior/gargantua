package block

import (
	"math/big"

	"github.com/EclesioMeloJunior/gargantua/keystore"
)

type Transaction struct {
	From keystore.Address
	To   keystore.Address

	Input  []*Transaction
	Output []*Transaction

	Fee        big.Int
	ReleasedAt int64 // timestamp
	Sig        []byte
}

// NewTx will generate a Transaction
func NewTx(i []*Transaction, o []*Transaction) *Transaction {
	return nil
}
