package networketh

import (
	"testing"
)
func TestConnection(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatal(err)
	} 
}