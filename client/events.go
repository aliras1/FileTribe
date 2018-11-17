package client

import (

	"github.com/golang/glog"

	"ipfs-share/eth"
	"bytes"
	. "ipfs-share/collections"
	"github.com/pkg/errors"
)

func (ctx *UserContext) HandleDebugEvents(ch chan *eth.EthDebug) {
	glog.Info("hadnling debug events...")

	for debug := range ch {
		glog.Infof("Eth Debug code: %v", debug.Msg.String())
	}
}


func (ctx *UserContext) HandleGroupInvitationEvents(ch chan *eth.EthGroupInvitation) {
	glog.Info("groupInvitation handling...")

	for groupInvitation := range ch {
		go ctx.onGroupInvitation(groupInvitation)
	}
}

func (ctx *UserContext) HandleGroupUpdateIpfsEvents(ch chan *eth.EthGroupUpdateIpfsHash) {
	glog.Info("groupUpdateIpfs handling...")

	for updateIpfs := range ch {
		ctx.onGroupUpdateIpfs(updateIpfs)
	}
}

func (ctx *UserContext) HandleGroupRegisteredEvents(ch chan *eth.EthGroupRegistered) {
	glog.Info("group registered handling...")

	for groupRegistered := range ch {
		ctx.onGroupRegisteredCallback(groupRegistered)
	}
}

func (ctx *UserContext) onGroupUpdateIpfs(updateIpfs *eth.EthGroupUpdateIpfsHash) {
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

func (ctx *UserContext) onGroupInvitation(inv *eth.EthGroupInvitation) {
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

func (ctx *UserContext) HandleKeyDirtyEvents(ch chan *eth.EthKeyDirty) {
	glog.Info("keyDirty handling...")

	for keyDirty := range ch {
		go ctx.onKeyDirty(keyDirty)
	}
}

func (ctx *UserContext) onKeyDirty(keyDirty *eth.EthKeyDirty) error {
	id := NewBytesId(keyDirty.GroupId)

	groupCtxInt := ctx.groups.Get(id)
	if groupCtxInt == nil {
		return nil
	}

	groupCtx := groupCtxInt.(*GroupContext)
	if err := groupCtx.Update(); err != nil {
		return errors.Wrap(err, "could not update group context")
	}

	// left or kicked. maybe send a signal to somewhere?
	if !groupCtx.Group.IsMember(ctx.user.Address()) {
		ctx.groups.DeleteWithId(id)

		return nil
	}

	groupCtx.On
}