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
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/collections"
	"github.com/aliras1/FileTribe/eth/gen/Account"
	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	"github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// AddressBook is a cache for fellow users' P2P contact data
type AddressBook struct {
	accountToContactMap *collections.Map
	backend             bind.ContractBackend
	app                 *ethapp.FileTribeDApp
	ipfs                ipfs.IIpfs
}

// NewAddressBook creates a new AddressBook
func NewAddressBook(backend bind.ContractBackend, app *ethapp.FileTribeDApp, ipfs ipfs.IIpfs) *AddressBook {
	return &AddressBook{
		backend:             backend,
		app:                 app,
		ipfs:                ipfs,
		accountToContactMap: collections.NewConcurrentMap(),
	}
}

// GetFromAccountAddress tries to retrieve the P2P contact data of a user based on its FT account
func (ab *AddressBook) GetFromAccountAddress(accountAddress ethcommon.Address) (*Contact, error) {
	var contact *Contact

	cInt := ab.accountToContactMap.Get(accountAddress)
	if cInt == nil {
		_c, err := ab.getContactFromEth(accountAddress)
		if err != nil {
			return nil, errors.Wrap(err, "could not get contact")
		}

		contact = _c
	} else {
		contact = cInt.(*Contact)
	}

	return contact, nil
}

// GetFromOwnerAddress tries to retrieve the P2P contact data of a user based on its owner
func (ab *AddressBook) GetFromOwnerAddress(ownerAddress ethcommon.Address) (*Contact, error) {
	for kv := range ab.accountToContactMap.KVIterator() {
		if bytes.Equal(kv.Value.(*Contact).OwnerAddress.Bytes(), ownerAddress.Bytes()) {
			return kv.Value.(*Contact), nil
		}
	}

	accountAddress, err := ab.app.GetAccountOf(&bind.CallOpts{Pending: true}, ownerAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not get account of owner")
	}

	contact, err := ab.getContactFromEth(accountAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not get contact from account address")
	}

	return contact, nil
}

func (ab *AddressBook) getContactFromEth(accountAddress ethcommon.Address) (*Contact, error) {
	acc, err := Account.NewAccount(accountAddress, ab.backend)
	if err != nil {
		return nil, errors.Wrap(err, "could not create account instance")
	}

	name, err := acc.Name(&bind.CallOpts{Pending: true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get account name")
	}

	ipfsID, err := acc.IpfsId(&bind.CallOpts{Pending: true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get ipfs peer id")
	}

	owner, err := acc.Owner(&bind.CallOpts{Pending: true})
	if err != nil {
		return nil, errors.Wrap(err, "could not get owner")
	}

	contact := NewContact(owner, accountAddress, name, ipfsID, tribecrypto.AnonymPublicKey{}, ab.ipfs)

	ab.accountToContactMap.Put(accountAddress, contact)

	return contact, nil
}
