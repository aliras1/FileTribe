package client

import (
	"bytes"
	"strings"
	"testing"

	nw "ipfs-share/networketh"
)

func TestBoxing(t *testing.T) {
	username := "testuser"
	password := "password"

	user := NewUser(username, password)

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

func TestSigning(t *testing.T) {
	username1 := "testuser1"
	password1 := "password1"

	user1 := NewUser(username1, password1)

	message := "Hello friend!"
	signedMessage := user1.Signer.SigningKey.Sign([]byte(message))
	msg, ok := user1.Signer.VerifyKey.Verify(signedMessage)
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
	ipfsAddress := "2434hasdf439asdjhbvc234f"
	network, err := nw.NewTestNetwork()
	if err != nil {
		t.Fatal(err)
	}
	user, err := SignUp(username, password, ipfsAddress, network)
	if err != nil {
		t.Fatal(err)
	}

	registered, err := network.IsUserRegistered(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("user should be registered")
	}

	_, uName, bKey, vKey, ipfs, err := network.GetUser(user.ID)
	if strings.Compare(uName, user.Name) != 0 {
		t.Fatal("usernames do not match")
	}
	if !bytes.Equal(bKey[:], user.Boxer.PublicKey.Value[:]) {
		t.Fatal("boxing keys do not match")
	}
	if !bytes.Equal(vKey[:], user.Signer.VerifyKey[:]) {
		t.Fatal("verify keys do not match")
	}
	if strings.Compare(ipfs, ipfsAddress) != 0 {
		t.Fatal("ipfs addresses do not match")
	}
}

func TestSignIn(t *testing.T) {
	username := "testuser"
	password := "password"
	network, err := nw.NewTestNetwork()
	if err != nil {
		t.Fatal(err)
	}
	_, err = SignUp(username, password, "ipfs", network)
	if err != nil {
		t.Fatal(err)
	}

	user, err := SignIn(username, password, network)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
}
