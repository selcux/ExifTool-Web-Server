package main

import (
	"context"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func handleWebSockets(ws *websocket.Conn) {
	log.Println("start websocket")
	readChan := make(chan string)
	convertChan := make(chan string)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	eixfUtil := NewEixfUtil(ctx)

	go eixfUtil.readExiftool(readChan)
	go eixfUtil.convertToJson(readChan, convertChan)

	for {
		msg := <-convertChan
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			log.Println("Client stopped listening...")
			readChan = nil
			cancel()
			break
		}
	}

	log.Println("end websocket")
}

func tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Incoming request")
	websocket.Handler(handleWebSockets).ServeHTTP(w, req)
	log.Println("Finished sending response...")
}

func main() {
	http.HandleFunc("/tags", tagHandler)

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
