package server

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var logServer *log.Logger = log.New(os.Stdout, "[Server] ", log.Ltime)

type safeBool struct {
	v   bool
	mux sync.Mutex
}

var serverStatus safeBool

// the channel size is the capacity of the queue
var msgQueueCh chan string = make(chan string, 100)

func ConsumeQueue() []string {
	var items []string
	queueIsEmpty := false
	logServer.Println("Queue is being consumed")
	for i := 1; true; i++ {
		if queueIsEmpty {
			break
		}
		select {
		case msg := <-msgQueueCh:
			// logServer.Printf("Item %d: %s\n", i, msg)
			items = append(items, msg)
		default:
			logServer.Println("Queue is empty")
			queueIsEmpty = true
		}
	}
	return items
}

func startListening() {
	serverStatus.mux.Lock()
	logServer.Println("Started listening to messages")
	serverStatus.v = true
	serverStatus.mux.Unlock()
}

func stopListening() {
	serverStatus.mux.Lock()
	logServer.Println("Stopped listening to messages")
	serverStatus.v = false
	serverStatus.mux.Unlock()
}

func IsListening() bool {
	return serverStatus.v
}

func listen(workerId int, ch <-chan string) {
	logPrefix := fmt.Sprintf("[Server][Worker %d] ", workerId)
	logWorker := log.New(os.Stdout, logPrefix, log.Ltime)
	for IsListening() {
		select {
		case msg, isOpen := <-ch:
			if !isOpen {
				logWorker.Println("Client has closed the channel.")
				stopListening()
				return
			}
			logServer.Printf("Accepted message '%s'\n", msg)
			// addToQueue(msg)
			select {
			case msgQueueCh <- msg:
				logWorker.Printf("Added message '%s' to queue.\n", msg)
			default:
				logWorker.Println("Queue is full. Discarding message.")
				stopListening()
				return
			}
		default:
			waitDuration := time.Duration(rand.Intn(100))
			time.Sleep(waitDuration * time.Millisecond)
		}
	}
}

func StartServer(numWorkers int) (chan<- string, error) {
	if numWorkers <= 0 {
		logServer.Printf("# workers must be positive (received %d)\n", numWorkers)
		return nil, errors.New("Received non-positive # of workers")
	}
	logServer.Printf("Starting server with %d workers\n", numWorkers)
	ch := make(chan string, numWorkers)
	startListening()
	for i := 1; i <= numWorkers; i++ {
		go listen(i, ch)
	}
	return ch, nil
}
