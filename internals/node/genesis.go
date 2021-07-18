package node

type Genesis struct {
	Chain    string             `json:"chainName"`
	Balances map[string]float32 `json:"balances"`
}
