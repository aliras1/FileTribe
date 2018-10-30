package client

import (

	"github.com/golang/glog"

	"ipfs-share/eth"
	"bytes"
	. "ipfs-share/collections"
	"encoding/base64"
)

func HandleDebugEvents(ctx *UserContext) {
	glog.Info("hadnling debug events...")

	for debug := range ctx.Network.GetDebugChannel() {
		glog.Infof("Eth Debug address: %v", debug.Msg)
	}
}

func HandleGroupInvitationEvents(ctx *UserContext) {
	glog.Info("groupInvitation handling...")

	for groupInvitation := range ctx.Network.GetGroupInvitationChannel() {
		go processGroupInvitationEvent(groupInvitation, ctx)
	}
}

func HandleGroupUpdateIpfsEvents(ctx *UserContext) {
	glog.Info("groupUpdateIpfs handling...")

	for updateIpfs := range ctx.Network.GetGroupUpdateIpfsChannel() {
		processGroupUpdateIpfsEvent(updateIpfs, ctx)
	}
}

func HandleGroupRegisteredEvents(ctx *UserContext) {
	glog.Info("group registered handling...")

	for groupRegistered := range ctx.Network.GetGroupRegisteredChannel() {
		processGroupRegisteredEvent(groupRegistered, ctx)
	}
}

func processGroupUpdateIpfsEvent(updateIpfs *eth.EthGroupUpdateIpfsPath, ctx *UserContext) {
	glog.Info("got update ipfs event message")

	groupCtxInterface := ctx.Groups.Get(NewBytesId(updateIpfs.GroupId))
	if groupCtxInterface == nil {
		return
	}

	groupCtx := groupCtxInterface.(*GroupContext)

	encNewIpfsHash, err := base64.URLEncoding.DecodeString(updateIpfs.IpfsPath)
	if err != nil {
		glog.Errorf("could not base64 decode new encrypted ipfs path")
		return
	}
	newIpfsHash, ok := groupCtx.Group.Boxer.BoxOpen(encNewIpfsHash)
	if !ok {
		glog.Errorf("could not decrpyt new ipfs hash")
		return
	}

	if err := groupCtx.Repo.Update(string(newIpfsHash)); err != nil {
		glog.Errorf("could not update group %s's repo with ipfs hash %s: %s", groupCtx.Group.Id.ToString(), updateIpfs.IpfsPath, err)
	}
}

func processGroupRegisteredEvent(groupRegistered *eth.EthGroupRegistered, ctx *UserContext) {
	glog.Info("got group registered event message")

	id := NewBytesId(groupRegistered.Id)
	groupCtxInterface := ctx.Groups.Get(id)
	if groupCtxInterface == nil {
		return
	}

	glog.Infof("group '%s' registered", id.ToString())
}

func processGroupInvitationEvent(inv *eth.EthGroupInvitation, ctx *UserContext) {
	glog.Infof("%s: got a group invitation event into %v", ctx.User.Name, inv.GroupId)

	// if new member
	if bytes.Equal(inv.To.Bytes(), ctx.User.Address.Bytes()) {
		glog.Infof("%s: entering group %v", ctx.User.Name, inv.GroupId)

		if err := NewGroupFromId(inv.GroupId, ctx); err != nil {
			glog.Warningf("could not start getting the group key: %s", err)
		}
	}

	// everyone updates
	groupCtxInterface := ctx.Groups.Get(NewBytesId(inv.GroupId))
	if groupCtxInterface == nil {
		return
	}

	glog.Infof("%s: updating group %v", ctx.User.Name, inv.GroupId)

	groupCtx := groupCtxInterface.(*GroupContext)
	if groupCtx == nil {
		glog.Warningf("no group found with id %v", inv.GroupId)
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Warningf("could not update group %v", inv.GroupId)
	}
}
