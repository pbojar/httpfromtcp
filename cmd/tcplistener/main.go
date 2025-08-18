package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	lst, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening for TCP traffic: %s\n", err)
	}
	defer lst.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := lst.Accept()
		if err != nil {
			log.Fatalf("Error connecting: %s", err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

// <- chan string: receive-only channel of strings
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go sendToChannel(f, ch)
	return ch
}

func sendToChannel(f io.ReadCloser, ch chan<- string) {
	defer f.Close()
	defer close(ch)
	currentLine := ""
	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil {
			if currentLine != "" {
				ch <- currentLine
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err)
			return
		}
		parts := strings.Split(string(b[:n]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			ch <- currentLine + parts[i]
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
}
