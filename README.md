# gofel

**Gofel** is a highly efficient, lightweight RPC (Remote Procedure Call) server framework written in Go. It is designed to simplify the process of setting up real-time communication servers with a focus on flexibility and performance. Leveraging the power of Gorilla WebSocket, gofel offers a robust solution for developers aiming to build scalable WebSocket applications.

## Features

- **Easy Setup**: Get your server up and running with minimal configuration.
- **Customizable Configurations**: Tailor the server to meet your specific needs through simple yet powerful configuration options.
- **WebSocket Integration**: Utilizes the well-known Gorilla WebSocket library for dependable and high-performance WebSocket connections.
- **Concurrent Processing**: Built to handle multiple requests concurrently, maximizing the efficiency of your real-time communications.
- **Lightweight Framework**: Focus on what's important - your application's functionality, not the boilerplate code.

## Getting Started

To get started with gofel, follow these simple steps:

### Prerequisites

Ensure you have Go installed on your machine (Go 1.15 or later is recommended). You can download it from [Go Downloads](https://golang.org/dl/).

### Installation

To use gofel in your project, simply import it:

```go
import "github.com/DanielcoderX/gofel"
```

### Running the Server

Create a file named `example.go` with the following content to run the server:

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/DanielcoderX/gofel/pkg/config"
	"github.com/DanielcoderX/gofel/api"
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
```

## Contribution

Contributions are what make the open-source community such a powerful platform for learning, inspiring, and creating. Any contributions you make are **greatly appreciated**.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Gorilla WebSocket library
- All contributors and supporters of the gofel project

We hope you find this framework useful. Happy coding!
