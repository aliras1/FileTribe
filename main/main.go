// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	ipfs_share "github.com/aliras1/FileTribe/client"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
)

type Config struct {
	APIAddress                 string
	IpfsAPIAddress             string
	EthFullNodeAddress         string
	EthAccountMnemonic         string
	EthAccountPasswordFilePath string
	FileTribeDAppAddress       string
	LogLevel                   string
}

const configPath = "./config.json"

var client ipfs_share.IUserFacade
var ipfs ipfsapi.IIpfs
var ethNode chequebook.Backend
var auth *ipfs_share.Auth

func signUp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	username := params["username"]
	err := client.SignUp(username)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("could not sign up: %s: %s", username, err))
		return
	}
}

func signOut(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	client.SignOut()
	client = nil
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	var groupname string
	if err := json.NewDecoder(r.Body).Decode(&groupname); err != nil {
		errorHandler(w, r, "argument not found")
		return
	}

	if err := client.CreateGroup(groupname); err != nil {
		errorHandler(w, r, err.Error())
	}
}

func joinGroup(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	var groupAddressStr string
	if err := json.NewDecoder(r.Body).Decode(&groupAddressStr); err != nil {
		errorHandler(w, r, "argument not found")
		return
	}

	if err := client.AcceptInvitation(ethcommon.HexToAddress(groupAddressStr)); err != nil {
		errorHandler(w, r, err.Error())
	}
}

func invite(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])

	var member string
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		errorHandler(w, r, "could not decode group id")
		return
	}
	address := ethcommon.HexToAddress(member)
	glog.Infof("inviting address: %s", address.String())

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			if err := group.Invite(address, true); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not invite user: %s", err.Error()))
			}
			return
		}
	}

	errorHandler(w, r, "no group found")
}

func leave(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			if err := group.Leave(); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not leave group: %s", err.Error()))
			}
			return
		}
	}

	errorHandler(w, r, "no group found")
}

func groupRepoCommit(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "")
		return
	}

	var path string
	json.NewDecoder(r.Body).Decode(&path)

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			if err := group.CommitChanges(); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not add group file: %s", err.Error()))
			}
			return
		}
	}

	errorHandler(w, r, "no group found")
}

func groupRepoListFiles(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "usercontext is nil")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			list := group.ListFiles()

			if err := json.NewEncoder(w).Encode(list); err != nil {
				errorHandler(w, r, "could not encode file list")
			}
			return
		}
	}

	errorHandler(w, r, "no group found")
}

func groupListMembers(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			var list []string
			for _, address := range group.ListMembers() {
				list = append(list, address.String())
			}

			glog.Error(list)

			if err := json.NewEncoder(w).Encode(list); err != nil {
				glog.Error(err)
			}

			return
		}
	}

	errorHandler(w, r, "no group found")
}

func groupRepoGrantWriteAccess(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])
	file := params["file"]
	address := ethcommon.HexToAddress(params["member"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			if err := group.GrantWriteAccess(file, address); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not grant write access: %s", err))
			}

			return
		}
	}

	errorHandler(w, r, "no group found")
}

func groupRepoRevokeWriteAccess(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupAddress := ethcommon.HexToAddress(params["groupAddress"])
	file := params["file"]
	address := ethcommon.HexToAddress(params["member"])

	for _, group := range client.Groups() {
		if bytes.Equal(group.Address().Bytes(), groupAddress.Bytes()) {
			if err := group.RevokeWriteAccess(file, address); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not grant write access: %s", err))
			}

			return
		}
	}

	errorHandler(w, r, "no group found")
}

func lsGroups(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is nil")
		return
	}

	var list []string
	for _, group := range client.Groups() {
		list = append(list, group.Address().String())
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not encode groupCtx list"))
	}
}

func listTransactions(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is nil")
		return
	}

	txList, err := client.Transactions()
	if err != nil {
		errorHandler(w, r, fmt.Sprintf(fmt.Sprintf("could not get tx's: %s", err)))
	}

	// just send back tx hashes for now
	var txs []string
	for _, tx := range txList {
		txs = append(txs, tx.Hash().Hex())
	}

	if err := json.NewEncoder(w).Encode(txs); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not encode txMap list: %s", err))
	}
}


func errorHandler(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}

func startDaemon() {
	var config Config

	jsonBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file: %s", err))
	}

	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		panic(fmt.Sprintf("could not unmarshal config json: %s", err))
	}

	if err := flag.Set("stderrthreshold", config.LogLevel); err != nil {
		panic(fmt.Sprintf("could not set log level: %s", err))
	}

	auth, err := ipfs_share.NewAuth(config.EthAccountMnemonic)
	if err != nil {
		panic(fmt.Sprintf("could not load account key data: NewNetwork: %s", err))
	}

	ipfs = ipfsapi.NewIpfs(config.IpfsAPIAddress)

	ethNode, err := ethclient.Dial(config.EthFullNodeAddress)
	if err != nil {
		panic(fmt.Sprintf("could not connect to ethereum node: %s", err))
	}

	client, err = ipfs_share.NewUserContext(
		auth,
		ethNode,
		ethcommon.HexToAddress(config.FileTribeDAppAddress),
		ipfs,
		"2001",
	)
	if err != nil {
		panic(fmt.Sprintf("could not create user context: %s", err))
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup/{username}", signUp).Methods("POST")
	router.HandleFunc("/signout", signOut).Methods("GET")

	router.HandleFunc("/group/create", createGroup).Methods("POST")
	router.HandleFunc("/group/join", joinGroup).Methods("POST")
	router.HandleFunc("/group/invite/{groupAddress}", invite).Methods("POST")
	router.HandleFunc("/group/leave/{groupAddress}", leave).Methods("POST")
	router.HandleFunc("/group/ls/{groupAddress}", groupListMembers).Methods("GET")
	router.HandleFunc("/group/repo/commit/{groupAddress}", groupRepoCommit).Methods("POST")
	router.HandleFunc("/group/repo/ls/{groupAddress}", groupRepoListFiles).Methods("GET")
	router.HandleFunc("/group/repo/grant/{groupAddress}/{file}/{member}", groupRepoGrantWriteAccess).Methods("POST")
	router.HandleFunc("/group/repo/revoke/{groupAddress}/{file}/{member}", groupRepoRevokeWriteAccess).Methods("POST")

	router.HandleFunc("/ls/groups", lsGroups).Methods("GET")
	router.HandleFunc("/ls/tx", listTransactions).Methods("GET")

	glog.Infof("serving on: %s", config.APIAddress)

	glog.Fatal(http.ListenAndServe(config.APIAddress, router))
}

func usage() {
	fmt.Print(`FileTribe

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
    EthAccountMnemonic                           Path to an Ethereum account key file
    EthAccountPasswordFilePath                  Path to the password file of the corresponding Ethereum account
    FileTribeDAppAddress                        Address of the FileTribeDApp contract
    LogLevel {INFO|WARNING|ERROR}               Level of logs that will be printed to stdout                                   

OPTIONS:
  -h --help                                     Show this screen`)
}

func printHelpAndExit(message string) {
	fmt.Printf("%s\nUse 'filetribe --help' to learn more about its usage.\n", message)
	os.Exit(1)
}

func main() {
	fileTribeUrl := flag.String("a", "http://127.0.0.1:3333", "")

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printHelpAndExit("No command found")
	}

	var (
		request *http.Request
		err error
	)
	url := *fileTribeUrl
	command := args[0]
	args = args[1:]

	switch command {
	case "signup":
		if len(args) < 1 {
			printHelpAndExit("No username found")
		}
		url += "/" + command + "/" + args[0]
		request, err = http.NewRequest("POST", url, bytes.NewBuffer(nil))
		if err != nil {
			panic(fmt.Sprintf("Could not create http request: %s", err))
		}

	case "ls":
		if len(args) < 1 {
			printHelpAndExit("You must specify what to list {-g|-i|-tx} (groups, pending invitations, pending transactions)")
		}

		switch args[0] {
		case "-g":
			url += "/" + command + "/groups"

		case "-i":
			url += "/" + command + "/invs"

		case "-tx":
			url += "/" + command + "/tx"

		default:
			printHelpAndExit("Argument must be one of {-g|-i|-tx}")
		}

		request, err = http.NewRequest("GET", url, bytes.NewBuffer(nil))
		if err != nil {
			panic(fmt.Sprintf("Could not create http request: %s", err))
		}

	case "daemon":
		startDaemon()

	case "group":
		if len(args) < 1 {
			printHelpAndExit("No group sub-command found")
		}

		url += "/" + command
		subcommand := args[0]
		args = args[1:]

		switch subcommand {
		case "create":
			if len(args) < 1 {
				printHelpAndExit("No group name found")
			}

			url += "/" + subcommand
			request, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`"%s"`, args[0]))))
			if err != nil {
				panic(fmt.Sprintf("Could not create http request: %s", err))
			}
			request.Header.Set("Content-Type", "application/json")

		case "invite":
			if len(args) < 2 {
				printHelpAndExit("Not enough arguments")
			}

			url += "/" + subcommand + "/" + args[0]
			request, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`"%s"`, args[1]))))
			if err != nil {
				panic(fmt.Sprintf("Could not create http request: %s", err))
			}
			request.Header.Set("Content-Type", "application/json")

		case "leave":
			if len(args) < 1 {
				printHelpAndExit("No group argument found")
			}

			url += "/" + subcommand + "/" + args[0]
			request, err = http.NewRequest("POST", url, bytes.NewBuffer(nil))
			if err != nil {
				panic(fmt.Sprintf("Could not create http request: %s", err))
			}
			request.Header.Set("Content-Type", "application/json")

		case "ls":
			if len(args) < 1 {
				printHelpAndExit("No group argument found")
			}

			url += "/" + subcommand + "/" + args[0]
			request, err = http.NewRequest("GET", url, bytes.NewBuffer(nil))
			if err != nil {
				panic(fmt.Sprintf("Could not create http request: %s", err))
			}

		case "repo":
			if len(args) < 2 {
				printHelpAndExit("Not enough arguments")
			}

			url += "/" + subcommand + "/" + args[0]

			subSubCommand := args[0]
			args = args[1:]

			switch subSubCommand {
			case "ls":
				url += "/" + subSubCommand + "/" + args[0]
				request, err = http.NewRequest("GET", url, bytes.NewBuffer(nil))
				if err != nil {
					panic(fmt.Sprintf("Could not create http request: %s", err))
				}
				request.Header.Set("Content-Type", "application/json")

			case "commit":
				url += "/" + subSubCommand + "/" + args[0]
				request, err = http.NewRequest("POST", url, bytes.NewBuffer(nil))
				if err != nil {
					panic(fmt.Sprintf("Could not create http request: %s", err))
				}
				request.Header.Set("Content-Type", "application/json")

			case "grant":
				fallthrough

			case "revoke":
				if len(args) < 3 {
					printHelpAndExit("Not enough arguments")
				}

				url += "/" + subSubCommand + "/" + args[0] + "/" + args[1] + "/" + args[2]
				request, err = http.NewRequest("POST", url, bytes.NewBuffer(nil))
				if err != nil {
					panic(fmt.Sprintf("Could not create http request: %s", err))
				}
				request.Header.Set("Content-Type", "application/json")
			}
		}

	default:
		printHelpAndExit("Unknown command")
	}

	if !strings.EqualFold(command, "daemon") {
		fmt.Printf("requesting: %s\n", url)


		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			panic(fmt.Sprintf("Could not send http request: %s", err))
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("could not read http response body: %s", err))
		}

		fmt.Println(string(body))
	}
}
