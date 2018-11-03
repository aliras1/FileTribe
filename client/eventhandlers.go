package client

import (

	"github.com/golang/glog"

	"ipfs-share/eth"
	"bytes"
	. "ipfs-share/collections"
)

func HandleDebugEvents(ctx *UserContext) {
	glog.Info("hadnling debug events...")

	for debug := range ctx.network.GetDebugChannel() {
		glog.Infof("Eth Debug code: %v", debug.Msg.String())
	}
}

func HandleGroupInvitationEvents(ctx *UserContext) {
	glog.Info("groupInvitation handling...")

	for groupInvitation := range ctx.network.GetGroupInvitationChannel() {
		go processGroupInvitationEvent(groupInvitation, ctx)
	}
}

func HandleGroupUpdateIpfsEvents(ctx *UserContext) {
	glog.Info("groupUpdateIpfs handling...")

	for updateIpfs := range ctx.network.GetGroupUpdateIpfsChannel() {
		processGroupUpdateIpfsEvent(updateIpfs, ctx)
	}
}

func HandleGroupRegisteredEvents(ctx *UserContext) {
	glog.Info("group registered handling...")

	for groupRegistered := range ctx.network.GetGroupRegisteredChannel() {
		processGroupRegisteredEvent(groupRegistered, ctx)
	}
}

func processGroupUpdateIpfsEvent(updateIpfs *eth.EthGroupUpdateIpfsHash, ctx *UserContext) {
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

func processGroupRegisteredEvent(groupRegistered *eth.EthGroupRegistered, ctx *UserContext) {
	glog.Info("got group registered event message")

	id := NewBytesId(groupRegistered.Id)
	groupCtxInterface := ctx.groups.Get(id)
	if groupCtxInterface == nil {
		return
	}

	glog.Infof("group '%s' registered", id.ToString())
}

func processGroupInvitationEvent(inv *eth.EthGroupInvitation, ctx *UserContext) {
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
