package client

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bgun26/msgq/server"
)

var logClient *log.Logger = log.New(os.Stdout, "[Client] ", log.Ltime)

func sendMsg(ch chan<- string, msg string) {
	waitDuration := time.Duration(rand.Intn(500))
	time.Sleep(waitDuration * time.Millisecond)
	logClient.Printf("Sending message '%s'\n", msg)
	ch <- msg
}

func closeChannel(ch chan<- string) {
	waitDuration := time.Duration(rand.Intn(2000) + 500)
	time.Sleep(waitDuration * time.Millisecond)
	logClient.Println("Closing channel")
	close(ch)
}

func SendMessages(ch chan<- string, numMessages int) {
	logClient.Printf("Sending %d messages to the server\n", numMessages)
	for i := 1; i <= numMessages; i++ {
		msg := fmt.Sprintf("Message #%d", i)
		go sendMsg(ch, msg)
	}
}

func ReadMessages(ch chan<- string) {
	go closeChannel(ch)
	for {
		if server.IsListening() {
			time.Sleep(50 * time.Millisecond)
		} else {
			logClient.Println("Server is not listening to messages.")
			break
		}
	}
	logClient.Println("Consuming the queue")
	queueItems := server.ConsumeQueue()
	for i, item := range queueItems {
		logClient.Printf("Item %d: %s\n", i+1, item)
	}

}
