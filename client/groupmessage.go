package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type GroupMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Proposal struct {
	From     string   `json:"from"`
	CMD      string   `json:"cmd"`
	Args     []string `json:"args"`
	PrevHash [32]byte `json:"prev_hash"`
	NewHash  [32]byte `json:"new_hash"`
}

type Approval struct {
	From string   `json:"from"`
	Hash [32]byte `json:"hash"`
}

func VerifyApproval(signedApproval []byte, network *nw.Network) (*Approval, [64]byte, error) {
	var signature [64]byte
	var approval Approval
	if err := json.Unmarshal(signedApproval[64:], &approval); err != nil {
		return nil, signature, errors.New("unmarshal: " + err.Error())
	}
	verifyKey, err := network.GetUserSigningKey(approval.From)
	if err != nil {
		return nil, signature, errors.New("could not get verify key: " + err.Error())
	}
	_, ok := verifyKey.Open(nil, signedApproval)
	if !ok {
		return nil, signature, errors.New("invalid approval")
	}
	copy(signature[:], signedApproval[:64])
	return &approval, signature, nil
}

func ValidateApproval(psm *ipfs.PubsubMessage, hash [32]byte, groupVerifyKey crypto.PublicSigningKey, network *nw.Network) (SignedBy, error) {
	signedBy := SignedBy{}
	signedApproval, ok := psm.Verify(groupVerifyKey)
	if !ok {
		return signedBy, errors.New("invalid group pubsub msg")
	}
	approval, signature, err := VerifyApproval(signedApproval, network)
	if err != nil {
		return signedBy, err
	}
	if !bytes.Equal(hash[:], approval.Hash[:]) {
		return signedBy, errors.New("invalid approval hash")
	}
	return SignedBy{approval.From, signature}, nil
}

type SignedBy struct {
	User      string   `json:"user"`
	Signature [64]byte `json:"signature"`
}

type CommitMsg struct {
	Proposal Proposal   `json:"proposal"`
	SignedBy []SignedBy `json:"signed_by"`
}
