package crypto

import (
	"crypto/rand"
	"fmt"
	"testing"

)

func TestSecretBox(t *testing.T) {
	//boxer := &SecretBoxer{
	//	Key: [32]byte{1},
	//}
	//
	//sigKey, err := ethcrypto.GenerateKey()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//signer := &Signer{PrivateKey: sigKey}
	//verifyKey := VerifyKey(ethcrypto.CompressPubkey(&signer.PrivateKey.PublicKey))
	//
	//f, err := os.Open("./boxingkeypair.go")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//defer f.Close()
	//
	//encReader, err := boxer.Seal(f, signer)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//foutName := "test_out.txt"
	//fout, err := os.Create(foutName)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//if err := boxer.Open(encReader, fout, &verifyKey); err != nil {
	//	t.Fatal(err)
	//}
	//
	//if err := fout.Close(); err != nil {
	//	t.Fatal(err)
	//}
	//
	//data1, _ := ioutil.ReadFile("./boxingkeypair.go")
	//data2, _ := ioutil.ReadFile(foutName)
	//
	//if !bytes.Equal(data1, data2) {
	//	t.Fatal("not equal")
	//}
	//
	//os.Remove(foutName)
}

func TestChunkNonce(t *testing.T) {
	var nonce [24]byte

	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatal(err)
	}

	fmt.Println(nonce)

	nonce, err := updateNonce(nonce)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(nonce)
}
