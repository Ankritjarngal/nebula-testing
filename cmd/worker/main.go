package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aym-n/nebula/types"
)

func main() {
	masterAddress := "localhost:8080"
	conn, err := net.Dial("tcp", masterAddress)
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Successfully connected to master at %s\n", masterAddress)

	workerLoop(conn)
}

func workerLoop(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	for {
		req := types.TaskRequest{WorkerID: "worker-123"}
		if err := enc.Encode(types.Message{Type: types.RequestTask, Payload: encodeGob(req)}); err != nil {
			log.Println("Error encoding task request:", err)
			break
		}

		var msg types.Message
		if err := dec.Decode(&msg); err != nil {
			log.Println("Error decoding task message:", err)
			break
		}

		if msg.Type == types.TaskResponseMsg {
			var task types.TaskResponse
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&task); err != nil {
				log.Println("Error decoding task payload:", err)
				break
			}

			// PlaceHolder
			log.Printf("Worker processing task: %s", task.TaskID)
			time.Sleep(1 * time.Second)

			result := types.ResultSubmission{
				TaskID: task.TaskID,
				Result: []byte("processed!"),
			}
			if err := enc.Encode(types.Message{Type: types.SubmitResult, Payload: encodeGob(result)}); err != nil {
				log.Println("Error encoding result submission:", err)
				break
			}
		} else {
			log.Printf("Unexpected message type received: %d", msg.Type)
		}
	}
}

func encodeGob(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		log.Fatal("Failed to encode data:", err)
	}
	return buf.Bytes()
}
