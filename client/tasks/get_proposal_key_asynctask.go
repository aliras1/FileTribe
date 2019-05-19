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

package tasks

import (
	"bytes"
	"encoding/base64"
	"github.com/aliras1/FileTribe/asynctask"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/aliras1/FileTribe/client/communication"
	"github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type getProposalKeyAsyncTask struct {
	account     interfaces.Account
	group       interfaces.Group
	proposalKey []byte
	p2p         *communication.P2PManager
	sessions    []common.Session
	eventCh     chan *asynctask.Event
	keyCh       chan tribecrypto.SymmetricKey
	stopCh      chan struct{}
	status      asynctask.Status
	lock        sync.Mutex
}

func NewGetProposalKeyAsyncTask(
	account interfaces.Account,
	group interfaces.Group,
	proposalKey string,
	p2pManager *communication.P2PManager,
	eventCh chan *asynctask.Event,
) asynctask.AsyncTask {
	return &getProposalKeyAsyncTask{
		account:     account,
		group:       group,
		proposalKey: []byte(proposalKey),
		p2p:         p2pManager,
		eventCh:     eventCh,
		keyCh:       make(chan tribecrypto.SymmetricKey, 100),
		stopCh:      make(chan struct{}),
		status:      asynctask.Pending,
	}
}

func (task *getProposalKeyAsyncTask) Execute() {
	task.lock.Lock()
	defer task.lock.Unlock()

	if task.status == asynctask.Pending {
		task.status = asynctask.Running
		go task.execute()
	}
}

func (task *getProposalKeyAsyncTask) execute() {
	glog.Infof("%s: ===== Started GetProposalKey task ======", task.account.Name())

	go func() {
		for _, memberOwner := range task.group.MemberOwners() {
			if bytes.Equal(memberOwner.Bytes(), task.account.Owner().Bytes()) {
				continue
			}

			session, err := task.p2p.StartGetProposedGroupKeySession(
				task.group.Address(),
				task.proposalKey,
				memberOwner,
				task.account.ContractAddress(),
				task.keyCh)

			if err != nil {
				glog.Errorf("could not start get group key session: %s", err)
			}

			task.sessions = append(task.sessions, session)
		}
	}()

	for {
		select {
		case key := <- task.keyCh:
			glog.Infof("%s: --> got key in task, %v : %v", task.account.Name(), task.eventCh, key.Key)
			if err := task.isBoxerValid(key); err != nil {
				glog.Infof("%s: invalid key", task.account.Name())
				continue
			}

			task.eventCh <- &asynctask.Event{
				Sender:task,
				Args:interfaces.Proposal{
					EncIpfsHash:task.proposalKey,
					Boxer:key,
				},
			}
			glog.Infof("%s: task sent back", task.account.Name())

			glog.Infof("%s: getProposedKeyAsyncTask stopped, cleaning up", task.account.Name())
			task.cleanUp()
			glog.Infof("%s: getProposedKeyAsyncTask cleaned up", task.account.Name())
			return

		case <- task.stopCh:
			glog.Infof("%s: getProposedKeyAsyncTask stopped, cleaning up", task.account.Name())
			task.cleanUp()
			glog.Infof("%s: getProposedKeyAsyncTask cleaned up", task.account.Name())
			return

		case <-time.After(5 * time.Second):
			glog.Infof("%s: getProposedKeyAsyncTask timeout, cleaning up", task.account.Name())
			task.cleanUp()
			glog.Infof("%s: getProposedKeyAsyncTask cleaned up", task.account.Name())
			return
		}
	}
}

func (task *getProposalKeyAsyncTask) isBoxerValid(proposedBoxer tribecrypto.SymmetricKey) error {
	if tribecrypto.BoxerIsNull(proposedBoxer) {
		return errors.New("boxer is null")
	}

	encIpfsHash, err := base64.StdEncoding.DecodeString(string(task.proposalKey))
	if err != nil {
		return errors.Wrap(err, "could not base64 decode enc ipfs hash")
	}

	_, ok := proposedBoxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encIpfsHash")
	}

	return nil
}

func (task *getProposalKeyAsyncTask) cleanUp() {
	if task.status == asynctask.Finished {
		return
	}

	for _, session := range task.sessions {
		session.Abort()
	}

	close(task.stopCh)
	close(task.keyCh)

	task.status = asynctask.Finished
}

func (task *getProposalKeyAsyncTask) Cancel() {
	if task.status == asynctask.Finished {
		return
	}

	glog.Infof("%s: Stopping getProposedKeyAsyncTask", task.account.Name())

	task.stopCh <- struct{}{}
	glog.Infof("%s: stop signal sent to getProposedKeyAsyncTask", task.account.Name())
}