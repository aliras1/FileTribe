package client

import (
	"strings"
	"testing"

	"bytes"
	nw "ipfs-share/network"
)

func TestBoxing(t *testing.T) {
	username1 := "testuser1"
	password1 := "password1"
	username2 := "testuser2"
	password2 := "password2"

	user1 := NewUser(username1, password1)
	user2 := NewUser(username2, password2)

	message := "Hello friend!"
	encMsg := user1.Boxer.BoxSeal([]byte(message), &user2.Boxer.PublicKey)
	var nonce [24]byte
	copy(nonce[:], encMsg[:24])
	plain, success := user2.Boxer.BoxOpen(encMsg, &user1.Boxer.PublicKey)
	if !success {
		t.Fatal("could not decrypt message")
	}
	if strings.Compare(string(plain), message) != 0 {
		t.Fatal("the original and the decrypted messages are not the same")
	}
}

func TestSiging(t *testing.T) {
	username1 := "testuser1"
	password1 := "password1"

	user1 := NewUser(username1, password1)

	message := "Hello friend!"
	signedMessage := user1.Signer.SecretKey.Sign(nil, []byte(message))
	msg, ok := user1.Signer.PublicKey.Open(nil, signedMessage)
	if !ok {
		t.Fatal("failed to verify")
	}
	if !bytes.Equal(msg, []byte(message)) {
		t.Fatal("messages do not match")
	}
}

func TestUserDataOnServer(t *testing.T) {
	username := "testuser"
	password := "password"
	network := nw.Network{"http://0.0.0.0:6000"}
	user, err := SignUp(username, password, "ipfs", &network)
	if err != nil {
		t.Fatal(err)
	}
	// Public key hash
	publicKeyHash, err := network.GetUserPublicKeyHash(username)
	if err != nil {
		t.Fatal(err)
	}
	if !publicKeyHash.Equals(&user.PublicKeyHash) {
		t.Fatal("the public key hashes do not match")
	}
	// Public signing key
	publicSigningKey, err := network.GetUserVerifyKey(username)
	if err != nil {
		t.Fatal(err)
	}
	if !publicSigningKey.Equals(&user.Signer.PublicKey) {
		t.Fatal("the public signing keys do not match")
	}
	// Public boxing key
	publicBoxingKey, err := network.GetUserBoxingKey(username)
	if err != nil {
		t.Fatal(err)
	}
	if !publicBoxingKey.Equals(&user.Boxer.PublicKey) {
		t.Fatal("the public boxing keys do not match")
	}
	// Ipfs address
	ipfsAddr, err := network.GetUserIPFSAddr(username)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(ipfsAddr, "ipfs") != 0 {
		t.Fatal("the public boxing keys do not match")
	}
}

func TestSignIn(t *testing.T) {
	username := "testuser3"
	password := "password3"
	network := nw.Network{"http://0.0.0.0:6000"}
	_, err := SignUp(username, password, "ipfs", &network)
	if err != nil {
		t.Fatal(err)
	}

	user_in, err := SignIn(username, password, &network)
	if err != nil {
		t.Fatal(err)
	}
	if user_in == nil {
		t.Fatal("user is nil")
	}
}
