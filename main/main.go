package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"net/http"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	ipfsapi "github.com/ipfs/go-ipfs-api"

	"ipfs-share/client"
	nw "ipfs-share/network"
)

const (
	host           = "0.0.0.0"
	port           = "3333"
	ipfsAPIAddress = "http://127.0.0.1:5001"
	// ethWSAddress   = "ws://127.0.0.1:8001"
	ethWSAddress   = "ws://172.18.0.2:8001"
	// ethKeyPath = "../misc/eth/ethkeystore/UTC--2018-05-19T19-06-19.498239404Z--c4f45f1822b614116ea5b68d4020f3ae1a0179e5"
	contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"
)

var userContext *client.UserContext
var network *nw.Network
var ipfs *ipfsapi.Shell

func stringToAddress(addressString string) ethcommon.Address {
	addressBytes := ethcommon.FromHex(addressString)
	address := ethcommon.BytesToAddress(addressBytes)

	return address
}

func addressToFriend(address ethcommon.Address) *client.Friend {
	for _, friend := range userContext.Friends {
		if bytes.Equal(friend.Contact.Address.Bytes(), address.Bytes()) {
			return friend
		}
	}

	return nil
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
		glog.Fatalf("could not connect to ethereum network: %s", err)
	}

	userContext, err = client.NewUserContextFromSignUp(
		params["username"],
		params["password"],
		ethKeyPath,
		params["username"], // data directory
		network,
		ipfs,
	)
	if err != nil {
		glog.Errorf("could not signup: %s", err)
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
	json.NewDecoder(r.Body).Decode(&ethKeyPath)
	
	network, err = nw.NewNetwork(ethWSAddress, ethKeyPath, contractAddress, "pwd")
	if err != nil {
		glog.Fatalf("could not connect to ethereum network: %s", err)
	}

	userContext, err = client.NewUserContextFromSignIn(
		params["username"],
		params["password"],
		ethKeyPath,
		params["username"], // data directory
		network,
		ipfs,
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

func addFriend(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params := mux.Vars(r)
	address := stringToAddress(params["address"])

	if err := userContext.AddFriend(address); err != nil {
		glog.Errorf("could not add friend: '%s': %s", address.String(), err)
	}
}

func addFile(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var path string
	json.NewDecoder(r.Body).Decode(&path)

	id, err := userContext.AddFile(path)
	if err != nil {
		glog.Errorf("could not add file: '%s': %s", path, err)
	}

	idURLEnc := base64.URLEncoding.EncodeToString(id[:])
	if err := json.NewEncoder(w).Encode(idURLEnc); err != nil {
		glog.Error(err)
	}
}

func sharePTP(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params := mux.Vars(r)
	address := stringToAddress(params["address"])
	friend := addressToFriend(address)

	if friend == nil {
		glog.Errorf("friend not found '%s'", address)
		return
	}

	glog.Info("Id: ", string(params["id"]))

	idBytes, err := base64.URLEncoding.DecodeString(params["id"])
	if err != nil {
		glog.Error(err)
		
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var id [32]byte
	copy(id[:], idBytes)

	
	if err := userContext.Repo[userContext.User.Address][id].Share(friend, userContext); err != nil {
		glog.Errorf("could not share file '%s' with '%s': %s", params["id"], address.String(), err)

		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func ls(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	list := userContext.List()
	if err := json.NewEncoder(w).Encode(list); err != nil {
		glog.Error(err)
	}
}

func getPendingFriends(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(userContext.PendingFriends); err != nil {
		glog.Error(err)
	}
}

func getWaitingFriends(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	list := []client.Friend{}
	for _, waiting := range userContext.WaitingFriends {
		list = append(list, *waiting.Friend)
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		glog.Error(err)
	}
}

func getFriends(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(userContext.Friends); err != nil {
		glog.Error(err)
	}
}

func confirmFriend(w http.ResponseWriter, r *http.Request) {
	if userContext == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params := mux.Vars(r)
	address := stringToAddress(params["address"])

	for _, waiting := range userContext.WaitingFriends {
		if bytes.Equal(waiting.Friend.Contact.Address.Bytes(), address.Bytes()) {
			if err := waiting.Confirm(userContext); err != nil {
				glog.Error(err)

				w.WriteHeader(http.StatusNotFound)
			}
			return
		}
	}

	glog.Warningf("no waiting friend found with '%s'", address.String())
	
	w.WriteHeader(http.StatusNotFound)
	return
}


func main() {
	flag.Parse()

	ipfs = ipfsapi.NewShell(ipfsAPIAddress)

	router := mux.NewRouter()
	
	
	router.HandleFunc("/signup/{username}/{password}", signUp).Methods("POST")
	router.HandleFunc("/signin/{username}/{password}", signIn).Methods("POST")
	router.HandleFunc("/signout", signOut).Methods("GET")
	
	router.HandleFunc("/get/friends", getFriends).Methods("GET")
	router.HandleFunc("/get/friends/pending", getPendingFriends).Methods("GET")
	router.HandleFunc("/get/friends/waiting", getWaitingFriends).Methods("GET")
	
	router.HandleFunc("/confirm/friend/{address}", confirmFriend).Methods("POST")

	router.HandleFunc("/add/friend/{address}", addFriend).Methods("POST")
	router.HandleFunc("/add/file", addFile).Methods("POST")
	
	router.HandleFunc("/share/ptp/{id}/{address}", sharePTP).Methods("POST")
	
	router.HandleFunc("/ls", ls).Methods("GET")

	
	glog.Fatal(http.ListenAndServe(host+":"+port, router))
}
