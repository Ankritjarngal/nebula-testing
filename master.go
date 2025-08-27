package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Master node listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleWorkerConnection(conn)
	}
}

func handleWorkerConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Worker connected from %s\n", conn.RemoteAddr())

	// Allocate jobs to workers
}