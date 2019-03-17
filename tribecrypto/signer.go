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
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
)


type Signer struct {
	PrivateKey *ecdsa.PrivateKey
}

func(s *Signer) Sign(digest []byte) ([]byte, error) {
	return ethcrypto.Sign(digest, s.PrivateKey)
}

func (s *Signer) Verify(digest, signature []byte) bool {
	pk := ethcrypto.CompressPubkey(&s.PrivateKey.PublicKey)
	return ethcrypto.VerifySignature(pk, digest, signature[:len(signature) - 1])
}

type VerifyKey []byte

func (vk *VerifyKey) Verify(digest, signature []byte) bool {
	addr, err := ethcrypto.SigToPub(digest, signature)
	if err != nil {
		glog.Errorf("error while ecrecover: %s", err)
		return false
	}

	glog.Infof("ecrecover addr: %s", ethcrypto.PubkeyToAddress(*addr).String())

	return ethcrypto.VerifySignature(*vk, digest, signature[:len(signature) - 1])
}