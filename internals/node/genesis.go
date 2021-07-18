package node

type Genesis struct {
	Chain    string            `json:"chainName"`
	Balances map[string]uint32 `json:"balances"`
}
