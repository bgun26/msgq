package main

import (
	"flag"
	"fmt"
	"msgq/client"
	"msgq/server"
	"os"
)

func main() {
	// Define CLI arguments
	numWorkersPtr := flag.Int("num-workers", 1, "Number of workers to start the server with")
	numMessagesPtr := flag.Int("num-messages", 5, "Number of messages to send to server")
	// Parse CLI arguments
	flag.Parse()
	if *numWorkersPtr <= 0 {
		fmt.Printf("# workers must be positive (received %d)", *numWorkersPtr)
		os.Exit(1)
	}
	if *numMessagesPtr <= 0 {
		fmt.Printf("# messages must be positive (received %d)", *numMessagesPtr)
		os.Exit(1)
	}
	// Start the application
	fmt.Println("Starting main")
	ch, err := server.StartServer(*numWorkersPtr)
	if err != nil {
		fmt.Println("Error while starting server")
		os.Exit(1)
	}
	client.SendMessages(ch, *numMessagesPtr)
	client.ReadMessages(ch)
	fmt.Println("Ending main")
}
