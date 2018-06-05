package networketh

import (
	"fmt"
	"strings"
	big "math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"ipfs-share/eth"
)

const key = `{"address":"c4f45f1822b614116ea5b68d4020f3ae1a0179e5","crypto":{"cipher":"aes-128-ctr","ciphertext":"c47565906c488c5122c805a31a3e241d0839cda984903ec28aa07c8892deb5b0","cipherparams":{"iv":"d7814d0dc15a383630c0439c6ad2fea8"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"78d74296f7796969b5764bcfda6cf1cd2cd5bfc423fc0897313b9d23e7e0f219"},"mac":"d852362f275a61fd32acdf040a136a08dc0dc25ab69ddc3d54404b17e9f85826"},"id":"ce2a2147-38d2-4d99-95c1-4968ff6b7a0e","version":3}`

func connect() error {

	conn, err := ethclient.Dial("http://127.0.0.1:8000")
	if err != nil {
		return fmt.Errorf("could not connect to ethereum node: networketh.connect: %s", err)
	}

	storage, err := eth.NewEth(common.HexToAddress("0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"), conn)
	if err != nil {
		return fmt.Errorf("could not instantiate contract: networketh.connect: %s", err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(key), "pwd")
	if err != nil {
		return err
	}
	
	tx, err := storage.Store(auth, big.NewInt(120))
	if err != nil {
		return err
	}

	fmt.Printf("Transfer pending: 0x%x\n", tx.Hash())

	return nil
}