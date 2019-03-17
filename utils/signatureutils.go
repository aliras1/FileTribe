package utils

import "github.com/pkg/errors"

func SigToRSV(sig []byte) (r [32]byte, s [32]byte, v uint8, err error) {
	if len(sig) != 65 {
		err = errors.New("signature must be of length 65")
		return
	}

	copy(r[:], sig[:32])
	copy(s[:], sig[32:64])
	v = sig[64]

	return
}