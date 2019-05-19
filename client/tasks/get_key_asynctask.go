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
	"github.com/aliras1/FileTribe/asynctask"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/aliras1/FileTribe/client/communication"
	"github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type getKeyAsyncTask struct {
	account interfaces.Account
	group   interfaces.Group
	p2p     *communication.P2PManager
	sessions []common.Session
	eventCh   chan *asynctask.Event
	keyCh     chan tribecrypto.SymmetricKey
	stopCh    chan struct{}
	status    asynctask.Status
	lock      sync.Mutex
}

func NewGetGroupKeyAsyncTask(account interfaces.Account, group interfaces.Group, p2pManager *communication.P2PManager, eventCh chan *asynctask.Event) asynctask.AsyncTask {
	return &getKeyAsyncTask{
		account:account,
		group:group,
		p2p:p2pManager,
		eventCh:eventCh,
		keyCh:make(chan tribecrypto.SymmetricKey),
		stopCh:make(chan struct{}),
	}
}

func (task *getKeyAsyncTask) Execute() {
	task.lock.Lock()
	defer task.lock.Unlock()

	if task.status == asynctask.Pending {
		task.status = asynctask.Running
		go task.execute()
	}
}

func (task *getKeyAsyncTask) execute() {
	glog.Infof("===== Started GetKey task ======")
	go func() {
		for _, memberOwner := range task.group.MemberOwners() {
			if bytes.Equal(memberOwner.Bytes(), task.account.Owner().Bytes()) {
				continue
			}

			session, err := task.p2p.StartGetGroupKeySession(
				task.group.Address(),
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
			if err := task.group.CheckBoxer(key); err != nil {
				glog.Errorf("invalid boxer: %s", err)
				continue
			}

			task.eventCh <- &asynctask.Event{Sender: task, Args:key}
			task.cleanUp()
			return

		case <- task.stopCh:
			task.cleanUp()
			glog.Infof("getKeyAsyncTask stopped")
			return

		case <-time.After(5 * time.Second):
			task.cleanUp()
			glog.Infof("getKeyAsyncTask timeout")
			return
		}
	}
}

func (task *getKeyAsyncTask) cleanUp() {
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

func (task *getKeyAsyncTask) Cancel() {
	if task.status == asynctask.Finished {
		return
	}

	glog.Infof("Stopping getKeyAsyncTask")

	task.stopCh <- struct{}{}
}