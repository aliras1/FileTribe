package client

import (

	"github.com/golang/glog"

	"ipfs-share/eth"
	"bytes"
	. "ipfs-share/collections"
)

func HandleDebugEvents(ch chan *eth.EthDebug) {
	glog.Info("hadnling debug events...")

	for debug := range ch {
		glog.Infof("Eth Debug code: %v", debug.Msg.String())
	}
}

type OnGroupInvitationCallback func(registered *eth.EthGroupInvitation)

func HandleGroupInvitationEvents(ch chan *eth.EthGroupInvitation, callback OnGroupInvitationCallback) {
	glog.Info("groupInvitation handling...")

	for groupInvitation := range ch {
		callback(groupInvitation)
	}
}

type OnGroupUpdateIpfsCallback func(registered *eth.EthGroupUpdateIpfsHash)

func HandleGroupUpdateIpfsEvents(ch chan *eth.EthGroupUpdateIpfsHash, callback OnGroupUpdateIpfsCallback) {
	glog.Info("groupUpdateIpfs handling...")

	for updateIpfs := range ch {
		callback(updateIpfs)
	}
}

type OnGroupRegisteredCallback func(registered *eth.EthGroupRegistered)

func HandleGroupRegisteredEvents(ch chan *eth.EthGroupRegistered, callback OnGroupRegisteredCallback) {
	glog.Info("group registered handling...")

	for groupRegistered := range ch {
		callback(groupRegistered)
	}
}

func (ctx *UserContext) onGroupUpdateIpfsCallback(updateIpfs *eth.EthGroupUpdateIpfsHash) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	glog.Info("got update ipfs event message")

	groupCtxInterface := ctx.groups.Get(NewBytesId(updateIpfs.GroupId))
	if groupCtxInterface == nil {
		return
	}

	groupCtx := groupCtxInterface.(*GroupContext)

	boxer := groupCtx.Group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(updateIpfs.IpfsHash)
	if !ok {
		glog.Errorf("could not decrpyt new ipfs hash")
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group context: %s", err)
	}

	if err := groupCtx.Repo.update(string(newIpfsHash)); err != nil {
		glog.Errorf("could not update group %s's repo with ipfs hash %s: %s", groupCtx.Group.Id().ToString(), updateIpfs.IpfsHash, err)
	}
}

func (ctx *UserContext) onGroupRegisteredCallback(groupRegistered *eth.EthGroupRegistered) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	glog.Info("got group registered event message")

	id := NewBytesId(groupRegistered.Id)
	groupCtxInterface := ctx.groups.Get(id)
	if groupCtxInterface == nil {
		return
	}

	glog.Infof("group '%s' registered", id.ToString())
}

func (ctx *UserContext) onGroupInvitationCallback(inv *eth.EthGroupInvitation) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	glog.Infof("%s: got a group invitation event into %v", ctx.user.Name, inv.GroupId)

	// if new member
	if bytes.Equal(inv.To.Bytes(), ctx.user.Address().Bytes()) {
		glog.Infof("%s: entering group %v", ctx.user.Name, inv.GroupId)

		if err := NewGroupFromId(inv.GroupId, ctx); err != nil {
			glog.Warningf("could not start getting the group key: %s", err)
		}
	}

	// everyone updates
	groupCtxInterface := ctx.groups.Get(NewBytesId(inv.GroupId))
	if groupCtxInterface == nil {
		return
	}

	glog.Infof("%s: updating group %v", ctx.user.Name, inv.GroupId)

	groupCtx := groupCtxInterface.(*GroupContext)
	if groupCtx == nil {
		glog.Warningf("no group found with id %v", inv.GroupId)
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Warningf("could not update group %s", err)
	}
}
