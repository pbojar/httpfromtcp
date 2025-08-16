package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	msg, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}

	dat := make([]byte, 8)
	for {
		_, err := msg.Read(dat)
		if err == io.EOF {
			fmt.Println("Reached end of file!")
			break
		}
		fmt.Printf("read: %s\n", dat)
	}
}
