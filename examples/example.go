package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"gofel/pkg/config"
	"gofel/api"
)

func main() {
	// Define custom configuration for the server
	// This is where you can set paths and ports differently from defaults.
	myconfig := config.Config{
		Path: "/here", // Specify the path where the server should run
		Port: "8088",  // Specify the port on which the server will listen
	}

	// Load configuration with override values from 'myconfig'
	// This step allows customization of server settings.
	cfg := config.LoadConfig(myconfig)

	// Create a new RPC server using the specified configuration
	// This sets up the server with custom settings provided in cfg.
	server := api.NewServer(cfg)
	

	// Register 'echo' function that sends back received messages
	// This function is a basic example of an RPC function handling messages.
	server.RegisterFunction("echo", func(conn *websocket.Conn, data interface{}) error {
		// Marshal the incoming data into JSON format.
		encodedData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to encode data: %s", err)
			return err
		}

		// Send the JSON back to the client as a text message.
		err = conn.WriteMessage(websocket.TextMessage, encodedData)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
			return err
		}

		return nil
	})

	// Start the server and listen for incoming connections
	// Any failure in starting the server will be logged as a fatal error.
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
