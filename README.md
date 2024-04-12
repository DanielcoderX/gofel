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
	"github.com/DanielcoderX/gofel/internal/rpc"
	"github.com/DanielcoderX/gofel/pkg/config"
)

func main() {
	// Custom configuration for the server
	myconfig := config.Config{
		Path: "/here",
		Port: "8088",
	}

	// Load configuration with override
	cfg := config.LoadConfig(myconfig)

	// Create a new RPC server with the specified configuration
	server, err := rpc.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Register 'echo' function that sends back received messages
	rpc.RegisterFunction("echo", func(conn *websocket.Conn, data interface{}) error {
		encodedData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to encode data: %s", err)
			return err
		}

		err = conn.WriteMessage(websocket.TextMessage, encodedData)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
			return err
		}

		return nil
	})

	// Start the server
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
