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

package client

import (
	"github.com/ethereum/go-ethereum/common"
	// "crypto/rand"
	"github.com/ugorji/go/codec"
	"strings"
	"testing"
)

func TestAsd(t *testing.T) {
	var (
		handler    codec.CborHandle
		out        []byte
		message    = "hello friend"
		messageDec string
	)

	cborEnc := codec.NewEncoderBytes(&out, &handler)

	if err := cborEnc.Encode(message); err != nil {
		t.Fatal(err)
	}

	cborDec := codec.NewDecoderBytes(out, &handler)

	if err := cborDec.Decode(&messageDec); err != nil {
		t.Fatal(err)
	}
}

func TestBoxing(t *testing.T) {
	username := "testuser"

	user, _ := NewAccount(username, common.Address{}, nil)

	message := "Hello friend!"
	boxer := user.Boxer()
	encMsg, err := boxer.Seal([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	plain, err := boxer.Open(encMsg)
	if err != nil {
		t.Fatalf("could not decrypt message: %s", err)
	}
	if strings.Compare(string(plain), message) != 0 {
		t.Fatal("the original and the decrypted messages are not the same")
	}
}
