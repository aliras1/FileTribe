package tribecrypto

import (
	"fmt"
	"encoding/hex"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestNewMerkleTree(t *testing.T) {
	leaves := [][]byte{[]byte{0x00}, []byte{0x01}, []byte{0x02}, []byte{0x03}, []byte{0x04}, []byte{0x05}}

	for i, leaf := range(leaves) {
		leaves[i] = ethcrypto.Keccak256(leaf)
	}

	tree := NewMerkleTree(leaves)
	
	leaf := leaves[1]
	proof, err := tree.Prove(leaf)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("leaf: 0x%s\n", hex.EncodeToString(leaf))
	fmt.Println("proof")
	for _, p := range(proof) {
		fmt.Printf("0x%s\n", hex.EncodeToString(p))
	}
	fmt.Printf("root: 0x%s\n", hex.EncodeToString(tree.root.value))
}