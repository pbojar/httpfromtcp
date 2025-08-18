package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = ":42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost"+port)
	if err != nil {
		log.Fatalf("Error resolving UDP address: %s\n", err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error dialing UDP connection: %s\n", err)
	}
	defer udpConn.Close()

	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := stdin.ReadString(byte('\n'))
		if err != nil {
			log.Printf("Error reading line: %s\n", err)
		}
		_, err = udpConn.Write([]byte(line))
		if err != nil {
			log.Printf("Error writing line: %s\n", err)
		}
	}
}
