package crypto

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	var skB [32]byte
	n, err := rand.Read(skB[:])
	if err != nil {
		t.Fatal(err)
	}
	var skBB [32]byte
	copy(skBB[:], skB[:])
	fmt.Println(skB)
	fmt.Printf("read: '%d'\n", n)

	_, sk := Ed25519KeyPair(&skBB)

	m := []byte("asdasd")
	fmt.Printf("m len: '%d'\n", len(m))

	mS := sk.Sign(m)
	fmt.Printf("signe len: '%d'", len(mS))
}
