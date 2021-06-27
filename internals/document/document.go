package document

type Document struct {
	URL  string `rlp:"ipfs_url"`
	Size uint64 `rlp:"size"`
}

func NewDocument(url string) *Document {
	return &Document{URL: url}
}
