package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func parse_args() []string {
	// Discarding programm name
	arguments := os.Args[1:]
	if len(arguments) != 3 {
		fmt.Println("[!] Invalid number of arguments!")
		fmt.Println("[+] Usage: ./server <local port> <remote host> <remote port>")
		os.Exit(1)
	}
	return arguments
}

func handle_client(conn net.Conn, remote_host string) {
	for {
		defer conn.Close()
		// Create connection to the remote host
		remote, err := net.Dial("tcp", remote_host)
		defer remote.Close()
		if err != nil {
			log.Fatal(err)
		}
		// Using io.Copy to copy data from server to client and vice versa
		go func() { io.Copy(conn, remote) }()
		go func() { io.Copy(remote, conn) }()

	}
}

func main() {
	args := parse_args()
	local_host := "0.0.0.0:" + args[0]
	remote_host := args[1] + ":" + args[2]
	fmt.Printf("[+] Starting server at: %s\n", local_host)
	// Starting listener
	server, err := net.Listen("tcp", local_host)
	defer server.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("[+] Connection from: %s -> %s\n", conn.LocalAddr().String(), remote_host)
		go handle_client(conn, remote_host)
	}
}
