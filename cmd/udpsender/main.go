package main

import (
	"net"
	"log"
	"os"
	"bufio"
	"fmt"
)

const CONNECTION_TYPE = "udp"
const LISTEN_ADDR = "127.0.0.1:42069"

func main() {
	addr, err := net.ResolveUDPAddr(CONNECTION_TYPE, LISTEN_ADDR)
	if err != nil {
		log.Fatalf("Error resolving %s address %s: %s\n", CONNECTION_TYPE, LISTEN_ADDR, err)
		os.Exit(1)
	}

	conn, err := net.DialUDP(CONNECTION_TYPE, nil, addr)
	if err != nil {
		log.Fatalf("Error on DialUPD %s address\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", LISTEN_ADDR)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
