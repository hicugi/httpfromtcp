package main

import (
	"os"
	"log"
	"io"
	"strings"
	"errors"
	"fmt"
)

const filePath = "./messages.txt"

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
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Coun't open %s: %s\n", filePath, err)
		return
	}

	linesChan := getLinesChannel(file)

	for line := range linesChan {
		fmt.Println("read:", line)
	}
}
