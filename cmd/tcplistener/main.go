package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pbojar/httpfromtcp/internal/request"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error getting request from reader: %s\n", err)
		}
		fmt.Print(req.RequestLine.String())
		fmt.Print(req.Headers.String())
		if len(req.Body) != 0 {
			fmt.Printf("Body:\n%s\n", string(req.Body))
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}
