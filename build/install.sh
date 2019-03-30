#!/usr/bin/env bash

type truffle >/dev/null 2>&1 || { echo >&2 "I require Truffle but it's not installed. Aborting."; exit 1; }

echo [*] Creating environment...
go version
echo GOROOT: ${GOROOT}
echo GOPATH: ${GOPATH}

ws=./build/go_workspace
mkdir -p ${ws}
cd ${ws}
export GOPATH=$PWD

root=./src/github.com/aliras1/
mkdir -p ${root}
cd ${root}
ln -s ../../../../../. FileTribe
cd ../../../../../

echo [*] Getting dependencies...

go get -u github.com/golang-collections/collections
go get -u github.com/sergi/go-diff/...
go get -u github.com/pkg/errors
go get -u github.com/golang/glog
go get -u github.com/whyrusleeping/tar-utils
go get -u golang.org/x/crypto/...
go get -u github.com/ethereum/go-ethereum
go get -u github.com/gorilla/mux
go get -u github.com/ipfs/go-ipfs-api
go get -u github.com/ugorji/go/codec
go get -u github.com/miguelmota/go-ethereum-hdwallet
go get -u github.com/tools/godep

CURRENT_DIR=$PWD
cd ${GOPATH}/src/github.com/ethereum/go-ethereum
${GOPATH}/bin/godep install ./cmd/abigen/
cd ${CURRENT_DIR}

echo [*] Generating abi APIs...

cd ./eth
./compile.sh

echo [*] Building FileTribe...

cd ../
mkdir ./build/bin

cd ./main
go build -o ../build/bin/filetribe main.go

echo [*] Creating symbolic link

cd ../
ln -s $PWD/build/bin/filetribe $HOME/.local/bin/filetribe

mkdir $HOME/.filetribe
cp ./config.json $HOME/.filetribe/