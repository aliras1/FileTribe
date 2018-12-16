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
	c := NewConcurrentMap()

	c.Put(&TestItem{&BytesId{[32]byte{1}}})
	c.Put(&TestItem{&BytesId{[32]byte{2}}})
	c.Put(&TestItem{&BytesId{[32]byte{3}}})
	c.Put(&TestItem{&BytesId{[32]byte{4}}})

	for item := range c.VIterator() {
		casted := item.(*TestItem)
		fmt.Println(casted.id)
	}

	c.Delete(&BytesId{[32]byte{3}})

	for item := range c.VIterator() {
		casted := item.(*TestItem)
		fmt.Println(casted.id)
	}
}
