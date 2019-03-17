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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/collections"
	"github.com/aliras1/FileTribe/eth/gen/Account"
	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	"github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type AddressBook struct {
	accToContactMap *collections.Map
	backend chequebook.Backend
	app *ethapp.FileTribeDApp
	ipfs ipfs.IIpfs
}

func NewAddressBook(backend chequebook.Backend, app *ethapp.FileTribeDApp, ipfs ipfs.IIpfs) *AddressBook {
	return &AddressBook{
		backend: 			backend,
		app:				app,
		ipfs:				ipfs,
		accToContactMap:	collections.NewConcurrentMap(),
	}
}

func (a *AddressBook) Get(accAddr ethcommon.Address) (*Contact, error) {
	var c *Contact

	cInt := a.accToContactMap.Get(accAddr)
	if cInt == nil {
		_c, err := a.getContactFromEth(accAddr)
		if err != nil {
			return nil, errors.Wrap(err, "could not get contact")
		}

		c = _c
	} else {
		c = cInt.(*Contact)
	}

	return c, nil
}

func (a *AddressBook) getContactFromEth(accAddr ethcommon.Address) (*Contact, error) {
	acc, err := Account.NewAccount(accAddr, a.backend)
	if err != nil {
		return nil, errors.Wrap(err, "could not create account instance")
	}

	name, err := acc.Name(&bind.CallOpts{Pending:true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get account name")
	}

	ipfsId, err := acc.IpfsId(&bind.CallOpts{Pending:true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get ipfs peer id")
	}

	owner, err := acc.Owner(&bind.CallOpts{Pending:true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get owner")
	}

	contact := NewContact(owner, accAddr, name, ipfsId, tribecrypto.AnonymPublicKey{}, a.ipfs)

	return contact, nil
}
