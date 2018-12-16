package client

import (
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	ethapp "ipfs-share/eth/gen/Dipfshare"
	ethgroup "ipfs-share/eth/gen/Group"
)

type Eth struct {
	Auth    *Auth
	App 	*ethapp.Dipfshare
	Backend chequebook.Backend
}

type GroupEth struct {
	*Eth
	Group *ethgroup.Group
}
