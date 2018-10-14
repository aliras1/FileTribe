package crypto

import (
	"testing"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"fmt"
	"encoding/hex"
)

func TestHashes(t *testing.T) {
	s1 := "hello"
	s2 := "bello"
	h := sha3.NewKeccak256()
	h.Write([]byte(s1))
	h.Write([]byte(s2))
	hash := h.Sum(nil)

	hexah := hex.EncodeToString(hash[:])

	fmt.Println(hexah)
}