package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message is a JSON-serializable structure that represents
// a request to be sent to the server.
type Message struct {
	Function string `json:"function"` // Name of the function to call on the server
	Data     string `json:"data"`      // Argument to pass to the function
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

	// Ensure that the connection is properly closed upon exit
	defer func() {
		// Cleanly close the connection by sending a close message
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		select {
		case <-time.After(time.Second):
		}
		c.Close()
	}()

	// Send an echo message
	echoMessage := Message{
		Function: "echo", // Name of the function to call on the server
		Data:     "Hello, world!", // Argument to pass to the function
	}
	echoMessageBytes, err := json.Marshal(echoMessage)
	if err != nil {
		log.Fatal("Error marshalling echo message:", err)
	}

	err = c.WriteMessage(websocket.TextMessage, echoMessageBytes)
	if err != nil {
		log.Fatal("write:", err)
	}

	// Read response
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}
	fmt.Printf("Received echo message: %s\n", message)
}

