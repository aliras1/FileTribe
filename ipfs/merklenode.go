package ipfs

type MerkleNode struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Size string `json:"size"`
	Data []byte `json:"data,omitempty"`
}
