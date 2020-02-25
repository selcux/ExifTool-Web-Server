package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

func main() {
	ws, err := websocket.Dial("ws://localhost:8000/tags", "", "http://localhost:8000/tags")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		fmt.Println(msg)
	}

	fmt.Println("Server finished request...")
}
