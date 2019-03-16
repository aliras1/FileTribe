package client

import (
	"github.com/ethereum/go-ethereum/contracts/chequebook"

	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
)

type Eth struct {
	Auth    *Auth
	App 	*ethapp.FileTribeDApp
	Backend chequebook.Backend
}

type GroupEth struct {
	*Eth
	Group *ethgroup.Group
}
