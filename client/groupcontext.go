package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"ipfs-share/client/communication/common"
	"path"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "ipfs-share/client/communication"
	sesscommon "ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	ipfsapi "ipfs-share/ipfs"
)

type IGroupFacade interface {
	Address() ethcommon.Address
	GrantWriteAccess(filePath string, user ethcommon.Address) error
	RevokeWriteAccess(filePath string, user ethcommon.Address) error
	CommitChanges() error
	Invite(user ethcommon.Address, hasInviteRigth bool) error
	Leave() error
	ListFiles() []string
	ListMembers() []ethcommon.Address
}

type GroupContext struct {
	account          interfaces.IAccount
	Group            interfaces.IGroup
	P2P              *com.P2PManager
	Repo             *fs.GroupRepo
	GroupConnection  *com.GroupConnection
	AddressBook      *common.AddressBook
	eth              *GroupEth
	Ipfs             ipfsapi.IIpfs
	Storage          *fs.Storage
	Transactions     *List
	broadcastChannel *ipfsapi.PubSubSubscription
	proposedKey      *crypto.SymmetricKey
	subs             *List
	lock             sync.Mutex
}

type GroupContextConfig struct {
	Group 			interfaces.IGroup
	Account 		interfaces.IAccount
	P2P 			*com.P2PManager
	AddressBook 	*common.AddressBook
	Eth 			*GroupEth
	Ipfs 			ipfsapi.IIpfs
	Storage 		*fs.Storage
	Transactions 	*List
}

func (groupCtx *GroupContext) Address() ethcommon.Address {
	return groupCtx.Group.Address()
}

func NewGroupContext(config *GroupContextConfig) (*GroupContext, error) {

	groupContext := &GroupContext{
		account:         config.Account,
		Group:           config.Group,
		P2P:             config.P2P,
		GroupConnection: nil,
		AddressBook:     config.AddressBook,
		eth:             config.Eth,
		Ipfs:            config.Ipfs,
		Storage:         config.Storage,
		Transactions:    config.Transactions,
		subs:            NewConcurrentList(),
	}

	repo, err := fs.NewGroupRepo(config.Group, config.Account.ContractAddress(), config.Storage, config.Ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not create group repo")
	}

	groupContext.Repo = repo
	groupContext.GroupConnection = com.NewGroupConnection(
		config.Group,
		repo,
		config.Account,
		config.AddressBook,
		onSessionClosed,
		config.P2P,
		config.Ipfs)

	go groupContext.HandleGroupInvitationSentEvents(config.Eth.Group)
	go groupContext.HandleGroupInvitationAcceptedEvents(config.Eth.Group)
	go groupContext.HandleNewConsensusEvents(config.Eth.Group)
	go groupContext.HandleIpfsHashChangedEvents(config.Eth.Group)
	go groupContext.HandleKeyDirtyEvents(config.Eth.Group)

	return groupContext, nil
}

func onSessionClosed(session sesscommon.ISession) {
	glog.Infof("session %d closed with error: %s", session.Id(), session.Error())
}

func (groupCtx *GroupContext) Update() error {
	contract := groupCtx.eth.Group

	name, err := contract.Name(&bind.CallOpts{Pending:true})
	if err != nil {
		return errors.Wrap(err, "could not get group name")
	}

	members, err := contract.Members(&bind.CallOpts{Pending:true})
	if err != nil {
		return errors.Wrap(err, "could not get group members")
	}

	encIpfsHash, err := contract.IpfsHash(&bind.CallOpts{Pending: true})
	if err != nil {
		return errors.Wrap(err, "could not get group ipfs hash")
	}


	if err := groupCtx.Group.Update(name, members, encIpfsHash); err != nil {
		return errors.Wrap(err, "could not update group")
	}

	if err := groupCtx.Group.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	if err := groupCtx.Repo.Update(groupCtx.Group.IpfsHash()); err != nil {
		return errors.Wrap(err, "could not update group repo")
	}

	return nil
}

func (groupCtx *GroupContext) Leave() error {
	tx, err := groupCtx.eth.Group.Leave(groupCtx.eth.Auth.TxOpts)
	if err != nil {
		return errors.Wrap(err, "could not send leave group tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
}

func (groupCtx *GroupContext) Stop() {
	groupCtx.GroupConnection.Kill()
}

func (groupCtx *GroupContext) CommitChanges() error {
	hash, err := groupCtx.Repo.CommitChanges(groupCtx.Group.Boxer())
	if err != nil {
		return errors.Wrap(err, "could commit group repo's changes")
	}

	key := groupCtx.Group.Boxer()
	encIpfsHash := key.BoxSeal([]byte(hash))

	tx, err := groupCtx.eth.Group.ChangeIpfsHash(groupCtx.eth.Auth.TxOpts, encIpfsHash)
	if err != nil {
		return errors.Wrap(err, "could not send change ipfs hash tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
}


func (groupCtx *GroupContext) Invite(newMember ethcommon.Address, hasInviteRight bool) error {
	glog.Infof("[*] Inviting account '%s' into group '%s'...\n", newMember, groupCtx.Group.Name)

	tx, err := groupCtx.eth.Group.Invite(groupCtx.eth.Auth.TxOpts, newMember)
	if err != nil {
		return errors.Wrap(err, "could not send invite account tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
}


func (groupCtx *GroupContext) Save() error {
	if err := groupCtx.Group.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	return nil
}


// Loads the locally available group meta data
func (groupCtx *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}

func (groupCtx *GroupContext) GrantWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not grant write access to non group members")
	}

	file := groupCtx.Repo.Get(NewStringId(path.Base(filePath)))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.account.ContractAddress()},
			groupCtx.Group.Address().String(),
			groupCtx.Storage,)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.GrantWriteAccess(groupCtx.account.ContractAddress(), user); err != nil {
		return errors.Wrap(err, "could not grant write access to account")
	}

	return nil
}

func (groupCtx *GroupContext) RevokeWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not revoke write access from non group members")
	}

	file := groupCtx.Repo.Get(NewStringId(path.Base(filePath)))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.account.ContractAddress()},
			groupCtx.Group.Address().String(),
			groupCtx.Storage,)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.RevokeWriteAccess(groupCtx.account.ContractAddress(), user); err != nil {
		return errors.Wrap(err, "could not revoke write access to account")
	}

	return nil
}

func (groupCtx *GroupContext) onKeyDirty() {
	glog.Info("KEY DIRTY")

	groupCtx.proposedKey = nil

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group: %s", err)
		return
	}

	newBoxer, err := crypto.NewSymmetricKey()
	if err != nil {
		glog.Errorf("could not create new group key: %s", err)
		return
	}

	groupCtx.proposedKey = newBoxer

	newIpfsHash, err := groupCtx.Repo.ReEncrypt(*newBoxer)
	if err != nil {
		glog.Errorf("could not re-encrypt group repo: %s", err)
		return
	}

	encNewIpfsHash := newBoxer.BoxSeal([]byte(newIpfsHash))

	tx, err := groupCtx.eth.Group.ChangeKey(groupCtx.eth.Auth.TxOpts, encNewIpfsHash)
	if err != nil {
		glog.Errorf("could not send change key tx: %s", err)

		groupCtx.proposedKey = nil

		return
	}

	simInt := interface{}(groupCtx.eth.Backend)
	x := simInt.(*backends.SimulatedBackend)
	x.Commit()

	groupCtx.Transactions.Add(tx)
}

func (groupCtx *GroupContext) ReEncrpyt() error {
	//hash, err := groupCtx.Repo.ReEncrypt(groupCtx.Group.Boxer())
	//if err != nil {
	//	return errors.Wrap(err, "could not re-encrypt group repo")
	//}
	//
	//if err := groupCtx.P2P.StartCommitSession(
	//	hash,
	//	groupCtx.account,
	//	groupCtx.Group,
	//	groupCtx.broadcast,
	//	groupCtx.OnCommitClientSuccess,
	//); err != nil {
	//	return errors.Wrap(err, "could not start new session")
	//}

	return errors.New("not implemented")
}

func (groupCtx *GroupContext) GetKey(encNewIpfsHash []byte) error {
	//newBoxer, ok := groupCtx.proposedKeys[encNewIpfsHashBase64]
	//
	//if ok {
	//	groupCtx.onGetKeySuccess(newBoxer)
	//} else {
	//	for _, member := range groupCtx.Group.Members() {
	//		if bytes.Equal(member.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
	//			continue
	//		}
	//
	//		c, err := groupCtx.AddressBook.Get(member)
	//		if err != nil {
	//			glog.Warningf("could not get contact for member: %s", member.String())
	//			continue
	//		}
	//
	//		if err := groupCtx.P2P.StartGetGroupKeySession(
	//			groupCtx.Group.Address(),
	//			c,
	//			groupCtx.account.ContractAddress(),
	//			func(cap *caps.GroupAccessCap) {
	//				groupCtx.onGetKeySuccess(cap.Boxer)
	//			},
	//		);	err != nil {
	//			glog.Errorf("could not start get group key session: %s", err)
	//		}
	//	}
	//}

	return nil
}

func (groupCtx *GroupContext) onGetKeySuccess(boxer crypto.SymmetricKey) {
	groupCtx.Group.SetBoxer(boxer)

	if err := groupCtx.Save(); err != nil {
		glog.Errorf("could not save new key: %s", err)
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group: %s", err)
	}
}

func (groupCtx *GroupContext) ListFiles() []string {
	var fileNames []string
	files := groupCtx.Repo.Files()

	for _, file := range files {
		fileNames = append(fileNames, file.Cap.FileName)
	}

	return fileNames
}

func (groupCtx *GroupContext) ListMembers() []ethcommon.Address {
	return groupCtx.Group.Members()
}

func (groupCtx *GroupContext) OnChangeGroupKeyClientSuccess(args []interface{}) {
	if len(args) < 1 {
		glog.Error("args should be min. of length 1")
	}

	encNewIpfsHash := args[0].([]byte)

	tx, err := groupCtx.eth.Group.ChangeKey(groupCtx.eth.Auth.TxOpts, encNewIpfsHash)
	if err != nil {
		glog.Errorf("could not send change group key transaction: %s", err)
		return
	}

	groupCtx.Transactions.Add(tx)
}

func (groupCtx *GroupContext) OnCommitClientSuccess(args []interface{}) {
	if len(args) < 1 {
		glog.Error("args should be of length 1")
	}

	encNewIpfsHash := args[0].([]byte)

	tx, err := groupCtx.eth.Group.ChangeIpfsHash(groupCtx.eth.Auth.TxOpts, encNewIpfsHash)
	if err != nil {
		glog.Errorf("could not send update group ipfs hash transaction: %s", err)
		return
	}

	groupCtx.Transactions.Add(tx)
}

func (groupCtx *GroupContext) broadcast(msg []byte) error {
	return groupCtx.GroupConnection.Broadcast(msg)
}

func (groupCtx *GroupContext) p2pBroadcast(msg []byte) error {
	for _, member := range groupCtx.Group.Members() {

		c, err := groupCtx.AddressBook.Get(member)
		if err != nil {
			glog.Warningf("could not get contact for member: %s", member)
		}

		go func() {
			if err := c.Send(msg); err != nil {
				glog.Errorf("error while sending p2p message: %s", err)
			}
		}()
	}

	return nil
}