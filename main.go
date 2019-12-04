package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bgun26/msgq/client"
	"github.com/bgun26/msgq/server"
)

const Version string = "0.1.1"

func usageFunc() {
	fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", Version)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

}

func main() {
	// Define CLI arguments
	numWorkersPtr := flag.Int("num-workers", 1, "Number of workers to start the server with")
	numMessagesPtr := flag.Int("num-messages", 5, "Number of messages to send to server")
	flag.Usage = usageFunc
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
