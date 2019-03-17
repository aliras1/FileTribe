# FileTribe

FileTribe is a blockchain-based decentralized file-sharing and editing system built on [Ethereum](https://www.ethereum.org/) and [IPFS](https://ipfs.io/).

##### :warning: WARNING :warning:
FileTribe is still in pre-alpha state. Do NOT use it for storing or sharing sensitive data. All branches are highly unstable.  

## Dependency

In order to use FileTribe, you will need a running IPFS daemon. To install IPFS, [download](https://dist.ipfs.io/#go-ipfs) it and then run
```
$ tar xvfz go-ipfs.tar.gz
$ cd go-ipfs
$ ./install.sh
```
For more information see [here](https://docs.ipfs.io/introduction/install/).

## Building the sources

To successfully build the FileTribe client application, you need Go (v1.10 <=) and Truffle. To install Truffle, run
`$ npm install -g truffle`. If both dependencies are installed, run

```
$ make all
```
This will download the Go dependencies, compile the solidity sources and create Go bindings to them. Note that you may have to use make in `sudo` mode since go-ethereum's abigen might fail when trying to resolve the dependencies of the generated Go files. 

## How to use

### Starting a single client instance

##### Deploy FileTribe

FileTribe is not deployed on any available public Ethereum networks currently since it is still under heavy developments.
If you want to try it you have to deploy it on a private network.
You can use your own or you can setup FileTribe's test network.
To do so you have to build geth in `FileTribe/build/go_workspace/src/github.com/ethereum/go-ethereum` by running `$ make geth` or use your own, already existing one.
Once you have geth, run the following snippet to bring up your private network: 

```
export FILE_TRIBE=/path/to/FileTribe
mkdir privnet
echo pwd > /tmp/password
geth --datadir ./privnet --keystore $FILE_TRIBE/ethkeystore/ init $FILE_TRIBE/ethgenesis/genesis.json
geth --datadir ./privnet --keystore $FILE_TRIBE/ethkeystore/ --networkid 15 --ws --wsaddr "0.0.0.0" --wsport "8001" --wsapi "db,eth,net,web3,personal,web3" --wsorigins "*" --rpc --rpcaddr "0.0.0.0" --rpcport "8000" --rpccorsdomain "*" --port "30304" --rpcapi "db,eth,net,web3,personal,web3" --nat "any" --nodiscover --password /password --unlock "c4f45f1822b614116ea5b68d4020f3ae1a0179e5" --mine --minerthreads=4
```

You can read more about building geth or creating your own network [here](https://github.com/ethereum/go-ethereum).
Now that you have a running network you can compile and deploy FileTribe on it. In an other terminal:

```
cd $FILE_TRIBE/eth
truffle compile --reset --all
truffle migrate --reset --network development
```

#### Start IPFS daemon

FileTribe uses IPFS as its data storage layer so the client needs an IPFS daemon which it can talk to. The clients communicate with each other through the IPFS daemon's built in _libp2p_ service which has to be enabled explicitly.   

```
ipfs init
ipfs daemon --enable-pubsub-experiment </dev/null &>/dev/null &
ipfs config --json Experimental.Libp2pStreamMounting true
```

#### Start the client application

If you have a running Ethereum network, on which the FileTribe contracts were deployed and an IPFS daemon you can start the client: `./filetribe <eth account key> [-ipfs=<addr>] [-eth=<addr>] [-p=<port>]`
 
`./filetribe /path/to/ethereum_account.key -ipfs=http://127.0.0.1:5001 -eth=ws://127.0.0.1:8001`

## Using Docker Compose

Alternatively, instead of manually setting up an Ethereum network, you can use FileTribe's prepared test environment which includes a full Ethereum node and three individual clients, all of them running in their separate containers. The following port mappings are available

* Alice: `3333` &rarr; `3333`
* Bob:  `3333`  &rarr; `3334`
* Charles:  `3333`  &rarr; `3335`

meaning you can access all three clients from your localhost. To use this simple environment type:

`docker-compose up`

## Execute commands on the client

To start using the application you have to start a client daemon process first. After that you can interact with that daemon.

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
    EthAccountKeyPath                           Path to an Ethereum account key file
    EthAccountPasswordFilePath                  Path to the password file of the corresponding Ethereum account
    FileTribeDAppAddress                        Address of the FileTribeDApp contract
    LogLevel {INFO|WARNING|ERROR}               Level of logs that will be printed to stdout                                   

OPTIONS:
  -h --help                                     Show this screen
```

## License

FileTribe is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also found in the `COPYING` file in the root of the repository.

#####Used libraries

| Library       | Author(s)        | License  |
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

The licenses of all the above mentioned libraries are included in the `COPYING.3RD-PARTY` in the root of the repository. 