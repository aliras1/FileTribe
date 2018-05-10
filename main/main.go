package main

import (
	"fmt"
	"net"
	"os"
	"log"

	"ipfs-share/client"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	ipfs, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		log.Printf("could not connect to ipfs daemon: %s", err)
	}
	network := &nw.Network{"http://172.18.0.2:6000"}
	var userContext *client.UserContext
	userContext = nil

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		fmt.Println("uzi elott")
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(&userContext, network, ipfs, conn)
	}
}

// Handles incoming requests.
func handleRequest(ctx **client.UserContext, netwrok *nw.Network, ipfs *ipfs.IPFS, conn net.Conn) {
	fmt.Println("handling uzi...")
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("error reading: handleRequest: %s", err)
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	cmd, err := client.NewCommand(string(buf))
	if err != nil {
		fmt.Printf("could not create new command handleRequest: %s", err)
		conn.Write([]byte(err.Error()))
		conn.Close()
		return
	}
	log.Println("[*] Executing api command...")
	uc, err := cmd.Execute(*ctx, netwrok, ipfs)
	log.Println("[*] Api command executed")
	if err != nil {
		log.Printf("asdasd")
		log.Printf("could not execute command: handleRequest: %s", err)
		conn.Write([]byte("could not execute command"))
		conn.Close()
		return
	}
	*ctx = uc
	// Close the connection when you're done with it.
	conn.Close()
}
