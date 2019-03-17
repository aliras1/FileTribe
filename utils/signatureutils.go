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