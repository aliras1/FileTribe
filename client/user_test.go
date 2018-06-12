package client

import (
	"bytes"
	// "crypto/rand"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	nw "ipfs-share/networketh"
)

func TestBoxing(t *testing.T) {
	username := "testuser"
	password := "password"

	user, _ := NewUser(username, password, "")

	message := "Hello friend!"
	encMsg, err := user.Boxer.Seal([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	plain, err := user.Boxer.Open(encMsg)
	if err != nil {
		t.Fatalf("could not decrypt message: %s", err)
	}
	if strings.Compare(string(plain), message) != 0 {
		t.Fatal("the original and the decrypted messages are not the same")
	}
}

func TestKeystore(t *testing.T) {
	// ks := keystore.NewKeyStore("../test/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// acc, err := ks.NewAccount("pwd")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fileName := path.Base(acc.URL.String())

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
	sig, err := user1.Signer.Sign(digest[:])
	if err != nil {
		t.Fatal(err)
	}


	pk := ethcrypto.CompressPubkey(&user1.Signer.PrivateKey.PublicKey)
	fmt.Println(pk)
	ok := ethcrypto.VerifySignature(pk, digest[:], sig[:64])
	if !ok {
		t.Fatal("failed to verify")
	}
}

func TestUserDataOnServer(t *testing.T) {
	username := "testuser"
	password := "password"
	ipfsAddress := "2434hasdf439asdjhbvc234f"

	dir := "../test/keystore"
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	keyAlice, ethKeyPath, err := nw.NewAccount(ks, dir, password)
	if err != nil {
		t.Fatal(err)
	}
	keyBob, _, err := nw.NewAccount(ks, dir, password)
	if err != nil {
		t.Fatal(err)
	}

	network, _, err := nw.NewTestNetwork(keyAlice, keyBob)
	if err != nil {
		t.Fatal(err)
	}

	user, err := SignUp(username, password, ethKeyPath, ipfsAddress, network)
	if err != nil {
		t.Fatal(err)
	}

	network.Simulator.Commit()

	registered, err := network.IsUserRegistered(user.Address)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("user should be registered")
	}

	uName, bKey, vKey, ipfs, err := network.GetUser(user.Address)
	if strings.Compare(uName, user.Name) != 0 {
		t.Fatal("usernames do not match")
	}
	if !bytes.Equal(bKey[:], user.Boxer.PublicKey.Value[:]) {
		t.Fatal("boxing keys do not match")
	}
	pk := ethcrypto.CompressPubkey(&user.Signer.PrivateKey.PublicKey)
	if !bytes.Equal(vKey, pk) {
		t.Fatal("verify keys do not match")
	}
	if strings.Compare(ipfs, ipfsAddress) != 0 {
		t.Fatal("ipfs addresses do not match")
	}

	// Test Sign in
	user, err = SignIn(username, password, ethKeyPath, network)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
}
