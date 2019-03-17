package client

import (
	// "crypto/rand"
	"github.com/ugorji/go/codec"
	"strings"
	"testing"
)


func TestAsd(t *testing.T) {
	var (
		handler codec.CborHandle
		out []byte
		message = "hello friend"
		messageDec string
	)

	cborEnc := codec.NewEncoderBytes(&out, &handler)

	if err := cborEnc.Encode(message); err != nil {
		t.Fatal(err)
	}

	cborDec := codec.NewDecoderBytes(out, &handler)

	if err := cborDec.Decode(&messageDec); err != nil {
		t.Fatal(err)
	}
}

func TestBoxing(t *testing.T) {
	username := "testuser"

	user, _ := NewAccount(username, nil)

	message := "Hello friend!"
	boxer := user.Boxer()
	encMsg, err := boxer.Seal([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	plain, err := boxer.Open(encMsg)
	if err != nil {
		t.Fatalf("could not decrypt message: %s", err)
	}
	if strings.Compare(string(plain), message) != 0 {
		t.Fatal("the original and the decrypted messages are not the same")
	}
}


