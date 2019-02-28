package tribecrypto

import (
	"fmt"
	// "ioioutil"
	"io"
	"crypto/rand"
	// "bufio"
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/sha3"


	"golang.org/x/crypto/nacl/secretbox"
)

const (
	chunkSize int = 104
	overheadSize int = 40
	maxUint uint64 = ^uint64(0) 
)

type FileBoxer struct {
	Key [32]byte  `json:"key"`
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

func (boxer *FileBoxer) SealAuth(reader io.Reader, signer *Signer) (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)
	var nonce [24]byte
	hasher := sha3.New256()

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

		if _, err := hasher.Write(enc); err != nil {
			return nil, fmt.Errorf("could not write to hasher: SecretBoxer.Seal: %s", err)
		}
		
		nonce, err = updateNonce(nonce)
		if err != nil {
			return nil, fmt.Errorf("could not update nonce: SecretBoxer.Seal: %s", err)
		}

		if _, err := buffer.Write(enc); err != nil {
			return nil, fmt.Errorf("could not write to buffer: SecretBoxer.Seal: %s", err)
		}
	}

	digest := hasher.Sum(nil)
	sig, err := signer.Sign(digest)
	if err != nil {
		return nil, fmt.Errorf("could not sign digest: SecretBoxer.Seal: %s", err)
	}

	finalBuffer := bytes.NewBuffer(sig[:64])
	if _, err := buffer.WriteTo(finalBuffer); err != nil {
		return nil, fmt.Errorf("could not create final buffer: SecretBoxer.Seal: %s", err)
	}

	return finalBuffer, nil
}

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

func (boxer *FileBoxer) OpenAuth(reader io.Reader, out io.Writer, verifyKey *VerifyKey) error {
	hasher := sha3.New256()

	var sig [64]byte
	if _, err := reader.Read(sig[:]); err != nil {
		return fmt.Errorf("could not read sig: SecretBoxer.Open: %s", err)
	}

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

		if _, err := hasher.Write(chunk[:n]); err != nil {
			return fmt.Errorf("could not write to hasher: SecretBoxer.Open: %s", err)
		}

		if _, err := out.Write(dec); err != nil {
			return fmt.Errorf("could not write to buffer: SecretBoxer.Open: %s", err)
		}	
	}

	digest := hasher.Sum(nil)

	if !verifyKey.Verify(digest, sig[:]) {
		return fmt.Errorf("could not verify: SecretBoxer.Open")
	}

	return nil
}

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
