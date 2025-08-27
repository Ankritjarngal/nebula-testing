package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	masterAddress := "localhost:8080"
	conn, err := net.Dial("tcp", masterAddress)
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Successfully connected to master at %s\n", masterAddress)
	
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Worker is still connected and ready for work.")

		// Recieve Job and process
	}
}