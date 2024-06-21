package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

// Message is a JSON-serializable structure that represents
// a request to be sent to the server.
type Message struct {
	Function string      `json:"function"` // Name of the function to call on the server
	Data     interface{} `json:"data"`     // Argument to pass to the function
}

func main() {
	// Adjust the URL to match your server's address and port
	url := "ws://localhost:8088/here" // Replace with your server's URL

	// Connect to the server
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Define a channel for clean shutdown
	done := make(chan struct{})
	defer close(done)

	// Define a function to send messages to the server
	sendMessage := func(funcName string, data interface{}, format string) {
		msg := Message{
			Function: funcName,
			Data:     data,
		}
		var msgBytes []byte
		if format == "msgpack" {
			msgBytes, err = msgpack.Marshal(msg)
		} else {
			msgBytes, err = json.Marshal(msg)
		}
		if err != nil {
			log.Printf("Error marshalling message for %s: %s", funcName, err)
			return
		}
		err = c.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("Failed to send message for %s: %s", funcName, err)
			return
		}
	}

	// Send an echo message
	sendMessage("echo", "Hello, world!", "json")
	// Read response
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Printf("read error: %s", err)
		return
	}
	log.Printf("Received echo message: %s\n", message)

	// Call hello function
	sendMessage("hello", nil, "json")
	// Read response
	_, message, err = c.ReadMessage()
	if err != nil {
		log.Printf("read error: %s", err)
		return
	}
	log.Printf("Received hello message: %s\n", message)

}
