package client

import (
	"fmt"
	"testing"
	"time"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

func Produce(c chan nw.Message) {
	for i := 0; i < 15; i++ {
		m := nw.Message{From: "from", Message: fmt.Sprintf("hello, friend!%d", i)}
		c <- m
	}
}

func TestMessageGetter(t *testing.T) {
	n := nw.Network{"http://0.0.0.0:6000"}
	i, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	NewUserContextFromUnamePassw("test_user", "pw", &n, i)

	n.SendMessage("from", "test_user", "hello friend!")
	n.SendMessage("from", "test_user", "hello friend, again!")
	fmt.Println("Sleeping...")
	time.Sleep(3 * time.Second)
	fmt.Println("End of test")
}
