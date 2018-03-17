package ipfs

import (
	"fmt"
	"testing"
)

func TestIPFS_AddDir(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	fmt.Println(err)
	ipfs.AddDir("/home/aliras/go/src/ipfs-share")
}
