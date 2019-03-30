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

package tribecrypto

import (
	"crypto/rand"
	"log"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/nacl/box"
)

func TestEncDec(t *testing.T) {
	pk2, sk2, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	boxer := AnonymBoxer{
		PublicKey:  AnonymPublicKey{pk2},
		PrivateKey: AnonymPrivateKey{sk2},
	}

	message := "Hello friend!"
	messageBox, err := boxer.PublicKey.Seal([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	m2, err := boxer.Open(messageBox)
	elapsed := time.Since(start)

	log.Println(elapsed)

	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(string(m2), message) != 0 {
		t.Fatal("messages do not match")
	}
}
