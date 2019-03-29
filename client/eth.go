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
	"github.com/ethereum/go-ethereum/contracts/chequebook"

	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
)

// Eth is a collection of those objects that are necessary for
// performing operations on the blockchain such as authentication
// data, DApp contract and a full ethereum node
type Eth struct {
	Auth    *Auth
	App 	*ethapp.FileTribeDApp
	Backend chequebook.Backend
}

// GroupEth stores a GroupContract and a pointer to all the
// general Eth data
type GroupEth struct {
	*Eth
	Group *ethgroup.Group
}
