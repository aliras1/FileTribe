package client

import (
	// "crypto/rand"
	"github.com/ugorji/go/codec"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

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
	password := "password"

	user, _ := NewUser(username, password, "")

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

func TestSigning(t *testing.T) {
	username1 := "testuser1"
	password1 := "password1"

	ks := keystore.NewKeyStore("../test/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(password1)
	if err != nil {
		t.Fatal(err)
	}
	fileName := "../test/keystore/" + path.Base(acc.URL.String())

	user1, err := NewUser(username1, password1, fileName)
	if err != nil {
		t.Fatal(err)
	}

	// message := "Hello friend!"
	digest := [32]byte{120}
	signer1 := user1.Signer()
	sig, err := signer1.Sign(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	pk := ethcrypto.CompressPubkey(&signer1.PrivateKey.PublicKey)
	fmt.Println(pk)
	ok := ethcrypto.VerifySignature(pk, digest[:], sig[:64])
	if !ok {
		t.Fatal("failed to verify")
	}
}

