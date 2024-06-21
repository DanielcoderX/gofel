package main

import (
	"context"
	"log"

	"github.com/DanielcoderX/gofel/api"
	"github.com/DanielcoderX/gofel/pkg/config"
	"github.com/DanielcoderX/gofel/pkg/wsconn"
)

func main() {
	// Define custom configuration for the server
	// This is where you can set paths and ports differently from defaults.
	myconfig := config.Config{
		Path:                   "/here", // Specify the path where the server should run
		Port:                   "8088",  // Specify the port on which the server will listen
		Verbose:                true,    // Enable verbose logging
		Format:                 "json",  // Use MSGPack as default format
		ConnectionPoolCapacity: 100,     // Default connection pool capacity
	}
	// Load configuration with override values from 'myconfig'
	// This step allows customization of server settings.
	cfg, err := config.LoadConfig(myconfig)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a new RPC server using the specified configuration
	// This sets up the server with custom settings provided in cfg.
	server := api.NewServer(cfg)

	// Register 'echo' function that sends back received messages
	// This function is a basic example of an RPC function handling messages.
	server.On("echo", func(conn *wsconn.ConnectionWrapper, data interface{}) {
		err := api.SendResponse(conn, data)
		if err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	})

	// This RPC function "hello" responds with a simple "Hello, world!"
	// message to all clients that call it.
	server.On("hello", func(conn *wsconn.ConnectionWrapper, data interface{}) {
		api.SendResponse(conn, "Hello Honey <3")
	})

	// Start the server and listen for incoming connections
	// Any failure in starting the server will be logged as a fatal error.
	ctx := context.Background() // Use background context as an example
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		server.Stop()
	}
}
