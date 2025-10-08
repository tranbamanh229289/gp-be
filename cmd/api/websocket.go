package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws() {
	http.HandleFunc("/ws", handleWebsocket)
	log.Println("Starting server at :8080 ... ")
	http.ListenAndServe(":8080", nil)

}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error", err)
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error", err)
			break;
		}

		log.Printf("Received: %s", message)
		response := []byte("Server received: " + string(message))
		err = conn.WriteMessage(messageType, response)
		if err != nil {
			log.Println("Write error", err)
			break
		}
	}

}