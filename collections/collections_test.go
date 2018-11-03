package collections

import (
	"testing"
	"fmt"
)

type TestItem struct {
	id IIdentifier
}

func (item *TestItem) Id() IIdentifier {
	return item.id
}

func TestConcurrentCollection(t *testing.T) {
	c := NewConcurrentCollection()

	c.Append(&TestItem{&BytesId{[32]byte{1}}})
	c.Append(&TestItem{&BytesId{[32]byte{2}}})
	c.Append(&TestItem{&BytesId{[32]byte{3}}})
	c.Append(&TestItem{&BytesId{[32]byte{4}}})

	for item := range c.Iterator() {
		casted := item.(*TestItem)
		fmt.Println(casted.id)
	}

	c.DeleteWithId(&BytesId{[32]byte{3}})

	for item := range c.Iterator() {
		casted := item.(*TestItem)
		fmt.Println(casted.id)
	}
}
