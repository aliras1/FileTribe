<p align="center">
<strong><em>WARNING:</em></strong>
FileTribe is still in pre-alpha state. Do <strong>NOT</strong> use it for storing or sharing sensitive data. All branches are highly unstable.
</p>
<hr/>

<p align="center">
  <img src="doc/logo.png" width="350">
</p>

## FileTribe

FileTribe is a blockchain-based decentralized file-sharing and editing system built on [Ethereum](https://www.ethereum.org/) and [IPFS](https://ipfs.io/).  

> Contract address of the development version on Ropsten: [0xec189ACCfbCAb44ED7b7665F293D45287120419b](https://ropsten.etherscan.io/address/0xec189ACCfbCAb44ED7b7665F293D45287120419b)

### Dependency

In order to use FileTribe, you will need a running [IPFS](https://ipfs.io/) daemon. To install [IPFS](https://ipfs.io/), [download](https://dist.ipfs.io/#go-ipfs) it and then run
```
$ tar xvfz go-ipfs.tar.gz
$ cd go-ipfs
$ ./install.sh
```
For more information see [here](https://docs.ipfs.io/introduction/install/).

### Building the sources

To successfully build the FileTribe client application, you need [Go](https://golang.org/) (v1.10 <=) and [Truffle](https://truffleframework.com/). To install [Truffle](https://truffleframework.com/), run
`$ npm install -g truffle`. For more information on installing [Truffle](https://truffleframework.com/), see [here](https://truffleframework.com/docs/truffle/getting-started/installation). If both dependencies are installed, run

```
$ make all
```
This will download the [Go](https://golang.org/) dependencies, compile the solidity sources and create [Go](https://golang.org/) bindings to them. Note that you may have to use make in `sudo` mode since go-ethereum's abigen might fail when trying to resolve the dependencies of the generated [Go](https://golang.org/) files. 

### How to use

An unsafe developer version of FileTribe is deployed on the **Ropsten** test network. The default `config.json` 
contains the necessary information to reach that application. 

1. To use FileTribe you need to edit `$HOME/.filetribe/config.json`
to include your mnemonic that generates your ethereum account. If you do not have any, you can generate one
using an [online mnemonic generator](https://iancoleman.io/bip39/) or by using [MetaMask](https://metamask.io/).
    ```json
    "EthAccountMnemonic": "your own long nice menmonic...",
    ```

2. To perform contract method calls you will need some `ether`. Use a faucet, like the one hosted by [MetaMask](https://metamask.io/)
to acquire ether.

    1. Navigate to [MetaMask's Test Ether Faucet](https://faucet.metamask.io/).

    2. Connect to the Ropsten Test Network using [MetaMask](https://metamask.io/).

    3. The faucet will link to your first account. Click `"Request 1 Ether From Faucet"` to submit your request.

    4. Within a short period of time, your account will be populated with the requested ether.

> Note, that FileTribe client uses [Infura](https://infura.io) to use the blockchain application by default. If you want it to use your own
specific [Ethereum](https://www.ethereum.org/) node or you deployed your own version of FileTribe and you want to use
that contract, you can set these in the `config.json` file.
 

##### Start IPFS daemon

FileTribe uses [IPFS](https://ipfs.io/) as its data storage layer so the client needs an [IPFS](https://ipfs.io/) daemon which it can talk to. The clients communicate with each other through the [IPFS](https://ipfs.io/) daemon's built in [_libp2p_](https://libp2p.io/) service which has to be enabled explicitly.   

```
ipfs init
ipfs daemon --enable-pubsub-experiment </dev/null &>/dev/null &
ipfs config --json Experimental.Libp2pStreamMounting true
```

##### Start the client application

You can configure the FileTribe daemon with `$HOME/.filetribe/config.json`. If you have a running [Ethereum](https://www.ethereum.org/) network, on which the FileTribe contracts were deployed and an [IPFS](https://ipfs.io/) daemon you can start the client:
 
```
$filetribe daemon
```

##### Execute commands on the client

Now that you have a running filetribe daemon you can start interacting with it. Since most of the operations you perform 
will result in a contract method call on the blockchain, these operations will not come into force in an instant.
You can check your pending transactions by executing `filetribe ls -tx` and lookup the results on [Etherscan's Ropsten part](https://ropsten.etherscan.io/). 

```
$filetribe --help
FileTribe

USAGE:
  filetribe <command> ...

COMMANDS: 
  BASIC COMMANDS:
    signup <username>                           Sign up to FileTribe    
    ls {-g|-i|-tx}                              List groups, pending invitations or pending Ethereum transactions
    daemon                                      Start a running client daemon process (configured from $HOME/.filetribe/config.json)                                                
    group                                       Interact with groups

  GROUP COMMANDS:
    create <groupname>                          Create a group
    invite <group address> <invitee address>    Invite a new member to the given group
    leave  <group address>                      Leave the given group
    ls <group address>                          List group members
    repo ...                                    Interact with the group repository

  REPO COMMANDS:
    ls <group address>                          List files
    commit <group address>                      Commit the pending changes in the repository
    grant <group address> <file> <member>       Grant write access for the given file to the given user
    revoke <group address> <file> <member>      Revoke write access for the given file to the given user

  CONFIG.JSON OPTIONS:
    APIAddress                                  Address on which the daemon will be listening    
    IpfsAPIAddress                              http address of a running IPFS daemon's API
    EthFullNodeAddress                          websocket address of an Ethereum full node
    EthAccountMnemonic                          Mnemonic that generates your Ethereum account
    EthAccountPasswordFilePath                  Path to the password file of the corresponding Ethereum account
    FileTribeDAppAddress                        Address of the FileTribeDApp contract
    LogLevel {INFO|WARNING|ERROR}               Level of logs that will be printed to stdout                                   

OPTIONS:
  -h --help                                     Show this screen
```

### License

FileTribe is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also found in the `COPYING` file in the root of the repository.

##### Used libraries

| Library       | Author(s)  | License  |
| ------------- |:----------:| --------:|
| [collections](https://github.com/golang-collections/collections)      | 2012 Caleb Doxsey | [MIT License](https://opensource.org/licenses/MIT) |
| [go-diff](https://github.com/sergi/go-diff)      | 2012-2016 The go-diff [Authors](https://github.com/sergi/go-diff/blob/master/AUTHORS)      |   [MIT License](https://opensource.org/licenses/MIT) |
| [errors](https://github.com/pkg/errors) | 2015, Dave Cheney <dave@cheney.net>      |    [BSD-2-Clause](https://opensource.org/licenses/BSD-2-Clause) |
| [glog](https://github.com/golang/glog) | Google      |    [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0) |
| [tar-utils](https://github.com/whyrusleeping/tar-utils) | 2016 Jeromy Johnson      |    [MIT License](https://opensource.org/licenses/MIT) |
| [go-ethereum](https://github.com/ethereum/go-ethereum) | 2014 The go-ethereum [Authors](https://github.com/ethereum/go-ethereum/blob/master/AUTHORS)      |    [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html) and [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html) |
| [mux](https://github.com/gorilla/mux) | 2012-2018 The Gorilla [Authors](https://github.com/gorilla/mux/blob/master/AUTHORS) |    [BSD 3-Clause "New" or "Revised" License](https://www.gnu.org/licenses/gpl-3.0.en.html) |
| [go-ipfs-api](https://github.com/ipfs/go-ipfs-api) | 2016 Jeromy Johnson |    [MIT License](https://opensource.org/licenses/MIT) |
| [go/codec](https://github.com/ugorji/go) | 2012-2015 Ugorji Nwoke |    [MIT License](https://opensource.org/licenses/MIT) |
| [go-ethereum-hdwallet](https://github.com/miguelmota/go-ethereum-hdwallet) | 2015 Miguel Mota |    [MIT License](https://opensource.org/licenses/MIT) |

The licenses of all the above mentioned libraries are included in the `COPYING.3RD-PARTY` in the root of the repository. 
