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

package common

import (
	"math"

	"github.com/aliras1/FileTribe/client/communication/common"
)


const EndOfSession = math.MaxUint8

type ISession interface {
	Id() uint32
	IsAlive() bool
	Abort()
	NextState(contact *common.Contact, data []byte)
	State() uint8
	Run()
	Error() error
}
