package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/aym-n/nebula/types"
)

type Task struct {
	ID    string
	Input []byte
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Master node listening on port 8080...")

	taskQueue := make(chan Task, 100)
	taskResults := make(chan types.ResultSubmission, 100)

	for i := 0; i < 100; i++ {
		taskID := fmt.Sprintf("task-%d", i)
		taskInput := []byte(fmt.Sprintf("input for %s", taskID))
		taskQueue <- Task{ID: taskID, Input: taskInput}
	}

	go func() {
		for result := range taskResults {
			log.Printf("Received result for task %s: %s", result.TaskID, string(result.Result))
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleWorkerConnection(conn, taskQueue, taskResults)
	}
}

func handleWorkerConnection(conn net.Conn, taskQueue <-chan Task, taskResults chan<- types.ResultSubmission) {
	defer conn.Close()
	fmt.Printf("Worker connected from %s\n", conn.RemoteAddr())

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	for {
		var msg types.Message

		if err := dec.Decode(&msg); err != nil {
			log.Printf("Error decoding message from worker %s: %v", conn.RemoteAddr(), err)
			return
		}

		switch msg.Type {
		case types.RequestTask:
			select {
			case task := <-taskQueue:
				taskResponse := types.TaskResponse{
					TaskID: task.ID,
					Input:  task.Input,
				}
				if err := enc.Encode(types.Message{Type: types.TaskResponseMsg, Payload: encodeGob(taskResponse)}); err != nil {
					log.Printf("Error encoding task response: %v", err)
				}
			default:
				log.Printf("No tasks available for worker %s", conn.RemoteAddr())
			}

		case types.SubmitResult:
			var result types.ResultSubmission
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&result); err != nil {
				log.Println("Error decoding result payload:", err)
				continue
			}
			taskResults <- result
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
