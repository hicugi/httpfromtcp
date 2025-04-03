package main

import "os"
import "io"
import "strings"
import "fmt"

func main() {
	file, _ := os.Open("./messages.txt")

	b8 := make([]byte, 8)
	line := ""

	for true {
		n, err := file.Read(b8)

		if err == io.EOF {
			break
		}

		str := string(b8[:n])
		idx := strings.Index(str, "\n")

		if idx == -1 {
			line += str
			continue
		}

		line += str[:idx]
		fmt.Printf("read: %s\n", line)

		line = str[idx+1:]
	}
}
