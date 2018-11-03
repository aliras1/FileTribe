package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"net/http"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	ipfs_share "ipfs-share/client"
	nw "ipfs-share/network"
	ipfsapi "ipfs-share/ipfs"
	"fmt"
	"ipfs-share/collections"
)

const (
	host           = "0.0.0.0"
	port           = "3333"
	ipfsAPIAddress = "http://127.0.0.1:5001"
	//ethWSAddress   = "ws://127.0.0.1:8001"
	ethWSAddress   = "ws://172.18.0.2:8001"
	// ethKeyPath = "../misc/eth/ethkeystore/UTC--2018-05-19T19-06-19.498239404Z--c4f45f1822b614116ea5b68d4020f3ae1a0179e5"
	contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"
)

var client ipfs_share.IUserFacade
var network nw.INetwork
var ipfs ipfsapi.IIpfs

func stringToAddress(addressString string) ethcommon.Address {
	addressBytes := ethcommon.FromHex(addressString)
	address := ethcommon.BytesToAddress(addressBytes)

	return address
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if client != nil {
		client.SignOut()
		client = nil
	}

	params := mux.Vars(r)
	var err error

	var ethKeyPath string
	json.NewDecoder(r.Body).Decode(&ethKeyPath)
	
	network, err = nw.NewNetwork(ethWSAddress, ethKeyPath, contractAddress, "pwd")
	if err != nil {
		errorHandler (w, r, fmt.Sprintf( "could not connect to ethereum network: %s", err))
		return
	}

	client, err = ipfs_share.NewUserContextFromSignUp(
		params["username"],
		params["password"],
		ethKeyPath,
		params["username"], // data directory
		network,
		ipfs,
		"2001",
	)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("could not signup: %s", err))
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	if client != nil {
		client.SignOut()
		client = nil
	}

	params := mux.Vars(r)
	var err error

	var ethKeyPath string
	if err := json.NewDecoder(r.Body).Decode(&ethKeyPath); err != nil {
		errorHandler(w, r, "could not decode ethkeypath")
		return
	}
	
	network, err = nw.NewNetwork(ethWSAddress, ethKeyPath, contractAddress, "pwd")
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("could not connect to ethereum network: %s", err))
		return
	}

	client, err = ipfs_share.NewUserContextFromSignIn(
		params["username"],
		params["password"],
		ethKeyPath,
		params["username"], // data directory
		network,
		ipfs,
		"2001",
	)
	if err != nil {
		glog.Error(err)
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

func invite(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	params := mux.Vars(r)
	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	var member string
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		errorHandler(w, r, "could not decode group id")
		return
	}
	address := ethcommon.HexToAddress(member)

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
			if err := group.Invite(address, true); err != nil {
				errorHandler(w, r, fmt.Sprintf("could not invite user: %s", err.Error()))
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
	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
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
	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
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
	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
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

	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	file := params["file"]
	address := ethcommon.HexToAddress(params["member"])

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
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

	groupIdArray, err := base64.URLEncoding.DecodeString(params["groupId"])
	if err != nil {
		errorHandler(w, r, "no group id found")
		return
	}
	var groupId [32]byte
	copy(groupId[:], groupIdArray)

	file := params["file"]
	address := ethcommon.HexToAddress(params["member"])

	for _, group := range client.Groups() {
		if group.Id().Equal(collections.NewBytesId(groupId)) {
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
		list = append(list, group.Id().ToString())
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

func main() {
	flag.Parse()

	ipfs = ipfsapi.NewIpfs(ipfsAPIAddress)

	router := mux.NewRouter()
	
	router.HandleFunc("/signup/{username}/{password}", signUp).Methods("POST")
	router.HandleFunc("/signin/{username}/{password}", signIn).Methods("POST")
	router.HandleFunc("/signout", signOut).Methods("GET")
	
	router.HandleFunc("/group/create", createGroup).Methods("POST")
	router.HandleFunc("/group/invite/{groupId}", invite).Methods("POST")
	router.HandleFunc("/group/ls/{groupId}", groupListMembers).Methods("GET")
	router.HandleFunc("/group/repo/commit/{groupId}", groupRepoCommit).Methods("POST")
	router.HandleFunc("/group/repo/ls/{groupId}", groupRepoListFiles).Methods("GET")
	router.HandleFunc("/group/repo/grant/{groupId}/{file}/{member}", groupRepoGrantWriteAccess).Methods("POST")
	router.HandleFunc("/group/repo/revoke/{groupId}/{file}/{member}", groupRepoRevokeWriteAccess).Methods("POST")

	router.HandleFunc("/ls/groups", lsGroups).Methods("GET")
	router.HandleFunc("/ls/tx", listTransactions).Methods("GET")
	
	glog.Fatal(http.ListenAndServe(host+":"+port, router))
}
