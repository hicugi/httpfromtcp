package main

import (
	"net"
	"log"
	"io"
	"strings"
	"errors"
	"fmt"
)

const CONNECTION_TYPE = "tcp"
const PORT = 42069

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		b8 := make([]byte, 8)
		currentLine := ""

		for {
			n, err := f.Read(b8)

			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(b8[:n])

			parts := strings.Split(str, "\n")
			lastIdx := len(parts)-1

			for i := 0; i < lastIdx; i++ {
				lines <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}

			currentLine += parts[lastIdx]
		}
	}()

	return lines
}

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

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println("read:", line)
		}

		// fmt.Println("- connect to", addr, "closed")
	}
}
