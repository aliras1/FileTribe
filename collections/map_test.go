// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

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
