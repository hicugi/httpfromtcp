package main

import (
	"net"
	"log"
	"fmt"
	"myhttpfromtcp/internal/request"
)

const CONNECTION_TYPE = "tcp"
const PORT = 42069

func main() {
	listener, err := net.Listen(CONNECTION_TYPE, fmt.Sprintf("%s%d", ":", PORT))
	if err != nil {
		log.Fatalf("Coun't create listener '%s' by port :%d\n", CONNECTION_TYPE, PORT)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("- connection error: %s\n", err.Error())
		}

		// addr := conn.RemoteAddr()
		// fmt.Println("- connection is live from", addr)

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("- error parsing request: %s\n", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}

		// fmt.Println("- connect to", addr, "closed")
	}
}
