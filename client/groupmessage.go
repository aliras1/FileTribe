package client

import (
	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type GroupMessage struct {
	From string `json:"from"`
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type HeartBeat struct {
	From string `json:"from"`
	Rand []byte `json:"rand"`
}

type Proposal struct {
	From     string   `json:"from"`
	CMD      string   `json:"cmd"`
	Args     []string `json:"args"`
	PrevHash [32]byte `json:"prev_hash"`
	NewHash  [32]byte `json:"new_hash"`
}

type Approval struct {
	From      string `json:"from"`
	Signature []byte `json:"signature"`
}

func (approval *Approval) Validate(rawTransaction []byte, groupSymKey crypto.SymmetricKey, network *nw.Network) error {
	// signed := append(approval.Signature, rawTransaction...)
	// verifyKey, err := network.GetUserVerifyKey(approval.From)
	// if err != nil {
	// 	return fmt.Errorf("could not get user verify key: ValidateApproval: %s", err)
	// }
	// _, ok := verifyKey.Open(nil, signed)
	// if !ok {
	// 	return fmt.Errorf("invalid approval: ValidateApproval")
	// }
	return nil
}
