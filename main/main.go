package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	ipfs_share "github.com/aliras1/FileTribe/client"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
)

var (
	host           = "0.0.0.0"
	port           = "3333"
	ipfsAPIAddress = "http://127.0.0.1:5001"
	//ethWSAddress   = "ws://127.0.0.1:8001"
	password = "pwd"
	ethWSAddress   = "ws://172.18.0.2:8001"
	ethKeyPath = "../misc/eth/ethkeystore/UTC--2018-05-19T19-06-19.498239404Z--c4f45f1822b614116ea5b68d4020f3ae1a0179e5"
	contractAddress = "0x41cf9ED28C99cC5eBd531bD1929a7E99c122fED8"
)

var client ipfs_share.IUserFacade
var ipfs ipfsapi.IIpfs
var ethNode chequebook.Backend
var auth *ipfs_share.Auth

func signUp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	username := params["username"]
	err := client.SignUp(username)
	glog.Flush()
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("could not sign up: %s: %s", username, err))
		return
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	//if client != nil {
	//	client.SignOut()
	//	client = nil
	//}
	//
	//params := mux.Vars(r)
	//var err error
	//
	//var ethKeyPath string
	//if err := json.NewDecoder(r.Body).Decode(&ethKeyPath); err != nil {
	//	errorHandler(w, r, "could not decode ethkeypath")
	//	return
	//}
	//
	//network, err = nw.NewNetwork(ethWSAddress, ethKeyPath, contractAddress, "pwd")
	//if err != nil {
	//	errorHandler(w, r, fmt.Sprintf("could not connect to ethereum network: %s", err))
	//	return
	//}
	//
	//client, err = ipfs_share.NewUserContextFromSignIn(
	//	params["username"],
	//	params["password"],
	//	ethKeyPath,
	//	params["username"], // data directory
	//	network,
	//	ipfs,
	//	"2001",
	//)
	//if err != nil {
	//	glog.Error(err)
	//}
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

	if err := json.NewEncoder(w).Encode(txList); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not encode txMap list: %s", err))
	}
}


func errorHandler(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}

func startDaemon(ethKeyPath string, ipfsApiAddress string, ethFullNodeAddress string, port string) {
	auth, err := ipfs_share.NewAuth(ethKeyPath, password)
	if err != nil {
		panic(fmt.Sprintf("could not load account key data: NewNetwork: %s", err))
	}

	ipfs = ipfsapi.NewIpfs(ipfsApiAddress)

	ethNode, err := ethclient.Dial(ethFullNodeAddress)
	if err != nil {
		if err != nil {
			panic(fmt.Sprintf("could not connect to ethereum node: %s", err))
		}
	}

	client, err = ipfs_share.NewUserContext(
		auth,
		ethNode,
		ethcommon.HexToAddress(contractAddress),
		ipfs,
		"2001",
	)
	if err != nil {
		panic(fmt.Sprintf("could not signup: %s", err))
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup/{username}", signUp).Methods("POST")
	router.HandleFunc("/signin/{username}", signIn).Methods("POST")
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

	glog.Fatal(http.ListenAndServe(host+":"+port, router))
}

func usage() {
	fmt.Print(`FileTribe

USAGE:
  filetribe <command> ...

COMMANDS: 
  BASIC COMMANDS:
    signup <username>                           Sign up to FileTribe    
    ls {-g|-i|-tx}                              List groups, pending invitations or pending Ethereum transactions
    daemon <eth account key> ...                Start a running client daemon process                                                
    group                                       Interact with groups

  GROUP COMMANDS:
    create <groupname>                          Create a group
    invite <group address> <invitee address>    Invite a new member to the given group
    leave  <group address>                      Leave the given group
    ls <group address>                          List group members
    repo <group address> ...                    Interact with the group repository

  REPO COMMANDS:
    ls                                          List files
    commit                                      Commit the pending changes in the repository
    grant <file> <member>                       Grant write access for the given file to the given user
    revoke <file> <member>                      Revoke write access for the given file to the given user

  DAEMON OPTIONS:
    -ipfs=<api address>                         http address of a running IPFS daemon's API
    -eth=<api address>                          websocket address of an Ethereum full node
    -p=<port>                                   Port number on which the daemon will be listening


OPTIONS:
  -h --help                                     Show this screen.
  -a=<address>                                  http address of a running client daemon`)
}

func printHelpAndExit(message string) {
	fmt.Printf("%s\nUse 'filetribe -h' to learn more about its usage.\n", message)
	os.Exit(1)
}

func main() {
	var err error

	fileTribeUrl := flag.String("a", "http://127.0.0.1:3333", "")

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printHelpAndExit("No command found")
	}

	request := *fileTribeUrl
	command := args[0]
	args = args[1:]

	switch command {
	case "signup":
		if len(args) < 2 {
			printHelpAndExit("No username found")
		}
		request += "/" + command + "/" + args[0]

	case "ls":
		if len(args) < 2 {
			printHelpAndExit("You must specify what to list {-g|-i|-tx} (groups, pending invitations, pending transactions)")
		}

		found := false
		split := strings.Split(args[0], "-")
		if len(split) < 2 {
			printHelpAndExit("Argument must be one of {-g|-i|-tx}")
		}

		for _, opt := range []string{"g", "i", "tx"} {
			if strings.EqualFold(opt, split[0]) {
				request += "/" + command + "/" + opt
				found = true
				break
			}
		}
		if !found {
			printHelpAndExit("Argument must be one of {-g|-i|-tx}")
		}

	case "daemon":
		fmt.Println(args)
		if len(args) < 1 {
			printHelpAndExit("No path to the Ethereum account key file found")
		}

		ethKeyPath = args[0]

		for _, arg := range args[1:] {
			split := strings.Split(arg, "=")
			if len(split) < 2 {
				printHelpAndExit(fmt.Sprintf("Invalid option: %s", arg))
			}

			opt, value := split[0], split[1]

			switch opt {
			case "-ipfs": ipfsAPIAddress = value
			case "-eth": ethWSAddress = value
			case "-p": port = value
			}

			startDaemon(ethKeyPath, ipfsAPIAddress, ethWSAddress, port)
		}

	case "group":
		if len(args) < 1 {
			printHelpAndExit("No group sub-command found")
		}

		request += "/" + command
		subcommand := args[0]
		args = args[1:]

		switch subcommand {
		case "create":
			if len(args) < 1 {
				printHelpAndExit("No group name found")
			}

			request += "/" + subcommand + "/" + args[0]

		case "invite":
			if len(args) < 2 {
				printHelpAndExit("Not enough arguments")
			}

			request += "/" + subcommand + "/" + args[0] + "/" + args[1]

		case "leave":
			fallthrough

		case "ls":
			if len(args) < 1 {
				printHelpAndExit("No group argument found")
			}

			request += "/" + subcommand + "/" + args[0]

		case "repo":
			if len(args) < 2 {
				printHelpAndExit("Not enough arguments")
			}

			request += "/" + subcommand + "/" + args[0]

			subSubCommand := args[0]
			args = args[1:]

			switch subSubCommand {
			case "ls":
				fallthrough

			case "commit":
				request += "/" + subSubCommand

			case "grant":
				fallthrough

			case "revoke":
				if len(args) < 2 {
					printHelpAndExit("Not enough arguments")
				}

				request += "/" + subSubCommand + "/" + args[0] + "/" + args[1]
			}
		}
	}

	resp, err := http.Get(*fileTribeUrl)
	if err != nil {
		fmt.Println(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))
	}
	defer resp.Body.Close()
}
