package crypto

import (
	"crypto/rand"
	"log"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/nacl/box"
)

func TestEncDec(t *testing.T) {
	pk2, sk2, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	boxer := AnonymBoxer{
		PublicKey:  AnonymPublicKey{pk2},
		PrivateKey: AnonymPrivateKey{sk2},
	}

	message := "Hello friend!"
	messageBox, err := boxer.PublicKey.Seal([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	m2, err := boxer.Open(messageBox)
	elapsed := time.Since(start)

	log.Println(elapsed)

	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(string(m2), message) != 0 {
		t.Fatal("messages do not match")
	}
}