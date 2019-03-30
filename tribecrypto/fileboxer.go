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
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// This should be forgotten and use a standard stream cipher

const (
	chunkSize    int    = 104
	overheadSize int    = 40
	maxUint      uint64 = ^uint64(0)
)

// FileBoxer is a stream cipher for encrypting huge files
type FileBoxer struct {
	Key [32]byte `json:"key"`
}

func updateNonce(base [24]byte) ([24]byte, error) {
	var int64Chunk1 uint64
	var int64Chunk2 uint64
	var int64Chunk3 uint64

	reader := bytes.NewReader(base[:])

	if err := binary.Read(reader, binary.BigEndian, &int64Chunk1); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert bin to int 1: chunkNonce %s", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &int64Chunk2); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert bin to int 2: chunkNonce %s", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &int64Chunk3); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert bin to int 3: chunkNonce %s", err)
	}

	if int64Chunk3 == maxUint {
		if int64Chunk2 == maxUint {
			int64Chunk1++
		}
		int64Chunk2++
	}
	int64Chunk3++

	buffer := bytes.NewBuffer(nil)

	if err := binary.Write(buffer, binary.BigEndian, &int64Chunk1); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert int to bin: chunkNonce %s", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &int64Chunk2); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert int to bin: chunkNonce %s", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &int64Chunk3); err != nil {
		return [24]byte{}, fmt.Errorf("could not convert int to bin: chunkNonce %s", err)
	}

	copy(base[:], buffer.Bytes())

	return base, nil
}

// Seal encrypts the contents of a file
func (boxer *FileBoxer) Seal(reader io.Reader) (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)
	var nonce [24]byte

	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("could not read random: SecretBoxer.Seal: %s", err)
	}

	for {
		var chunk [chunkSize - overheadSize]byte

		n, err := reader.Read(chunk[:])
		if n == 0 {
			break
		}

		enc := secretbox.Seal(nil, chunk[:n], &nonce, &boxer.Key)
		enc = append(nonce[:], enc...)

		nonce, err = updateNonce(nonce)
		if err != nil {
			return nil, fmt.Errorf("could not update nonce: SecretBoxer.Seal: %s", err)
		}

		if _, err := buffer.Write(enc); err != nil {
			return nil, fmt.Errorf("could not write to buffer: SecretBoxer.Seal: %s", err)
		}
	}

	return buffer, nil
}

// Open decrypts the provided cipher text
func (boxer *FileBoxer) Open(reader io.Reader, out io.Writer) error {
	for {
		var chunk [chunkSize]byte

		n, err := reader.Read(chunk[:])
		if n == 0 {
			break
		} else if err != nil {
			return fmt.Errorf("could not read from reader: SecretBoxer.Open: %s", err)
		}

		var nonce [24]byte
		copy(nonce[:], chunk[:24])

		dec, ok := secretbox.Open(nil, chunk[24:n], &nonce, &boxer.Key)
		if !ok {
			return fmt.Errorf("could not decrypt chunk: SecretBoxer.Open")
		}

		if _, err := out.Write(dec); err != nil {
			return fmt.Errorf("could not write to buffer: SecretBoxer.Open: %s", err)
		}
	}

	return nil
}
