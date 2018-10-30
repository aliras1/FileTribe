package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"net/http"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"ipfs-share/client"
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

var userContext *client.UserContext
var network nw.INetwork
var ipfs ipfsapi.IIpfs

func stringToAddress(addressString string) ethcommon.Address {
	addressBytes := ethcommon.FromHex(addressString)
	address := ethcommon.BytesToAddress(addressBytes)

	return address
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if userContext != nil {
		userContext.SignOut()
		userContext = nil
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

	userContext, err = client.NewUserContextFromSignUp(
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
	if userContext != nil {
		userContext.SignOut()
		userContext = nil
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

	userContext, err = client.NewUserContextFromSignIn(
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
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userContext.SignOut()
	userContext = nil
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		errorHandler(w, r, "user context is null")
		return
	}

	var groupname string
	if err := json.NewDecoder(r.Body).Decode(&groupname); err != nil {
		errorHandler(w, r, "argument not found")
		return
	}

	if err := userContext.CreateGroup(groupname); err != nil {
		errorHandler(w, r, err.Error())
	}
}

func invite(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
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

	groupCtxInterface := userContext.Groups.Get(collections.NewBytesId(groupId))
	if groupCtxInterface == nil {
		errorHandler(w, r, "no group found")
		return
	}
	groupCtx := groupCtxInterface.(*client.GroupContext)
	if err := groupCtx.Invite(address); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not invite user: %s", err.Error()))
	}
}

func addGroupFile(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
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


	groupCtxInterface := userContext.Groups.Get(collections.NewBytesId(groupId))
	if groupCtxInterface == nil {
		errorHandler(w, r, "no group found")
		return
	}
	groupCtx := groupCtxInterface.(*client.GroupContext)
	if err := groupCtx.AddFile(path); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not add group file: %s", err.Error()))
	}
}

func groupLs(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		errorHandler(w, r, "")
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

	groupCtxInterface := userContext.Groups.Get(collections.NewBytesId(groupId))
	if groupCtxInterface == nil {
		errorHandler(w, r, "no group found")
		return
	}

	groupCtx := groupCtxInterface.(*client.GroupContext)
	if err := json.NewEncoder(w).Encode(groupCtx.Repo.List()); err != nil {
		errorHandler(w, r, "could not encode file list")
	}
}


func ls(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		errorHandler(w, r, "")
		return
	}

	list := userContext.List()
	if err := json.NewEncoder(w).Encode(list); err != nil {
		glog.Error(err)
	}
}

func lsGroups(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		errorHandler(w, r, "user context is nil")
		return
	}

	var list []string
	for groupCtx := range userContext.Groups.Iterator() {
		list = append(list, groupCtx.(*client.GroupContext).Group.Name)
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		errorHandler(w, r, fmt.Sprintf("could not encode groupCtx list"))
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
	router.HandleFunc("/group/add/{groupId}", addGroupFile).Methods("POST")
	router.HandleFunc("/group/ls/{groupId}", groupLs).Methods("POST")
	
	router.HandleFunc("/ls", ls).Methods("GET")
	router.HandleFunc("/ls/groups", lsGroups).Methods("GET")
	
	glog.Fatal(http.ListenAndServe(host+":"+port, router))
}
