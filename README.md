# Gofel

**Gofel** is a high-performance, lightweight RPC (Remote Procedure Call) framework crafted in Go, aimed at streamlining the development of real-time communication servers. By harnessing the capabilities of the Gorilla WebSocket library, **Gofel** empowers developers to create scalable, efficient WebSocket applications with ease.

## Features

- **Easy Setup**: Quickly launch your server with minimal configuration.
- **Customizable Configurations**: Adapt the server settings to perfectly fit your application requirements.
- **WebSocket Integration**: Built on the reliable and high-performance Gorilla WebSocket library.
- **Concurrent Processing**: Designed for handling multiple requests simultaneously, enhancing real-time communication efficiency.
- **Lightweight Framework**: Focus more on your application's features and less on managing boilerplate code.

## Getting Started

Embark on your Gofel journey with these simple steps:

### Prerequisites

Make sure you have Go installed (version 1.15 or newer). Download it from [Go Downloads](https://golang.org/dl/).

### Installation

Incorporate Gofel into your project by adding the following import:

```go
import "github.com/DanielcoderX/gofel"
```

### Function Registration and Message Handling

When setting up your server with Gofel, you'll typically register functions that can handle incoming JSON requests in the following format:

```json
{
    "function": "echo",
    "data": "userdatatofunction"
}
```

Each registered function must be capable of parsing this structure to perform the required action. Here's an example of how to register such a function:

```go
server.RegisterFunction("echo", func(conn *websocket.Conn, data interface{}) error {
    var request struct {
        Function string `json:"function"`
        Data     string `json:"data"`
    }

    if err := json.Unmarshal(data.([]byte), &request); err != nil {
        log.Printf("Error parsing request: %s", err)
        return err
    }

    // Process the request based on 'Function' and 'Data'
    response, err := json.Marshal(map[string]interface{}{
        "result": "Received: " + request.Data,
    })
    if err != nil {
        log.Printf("Failed to encode response: %s", err)
        return err
    }

    err = conn.WriteMessage(websocket.TextMessage, response)
    if err != nil {
        log.Printf("Failed to send message: %s", err)
        return err
    }

    return nil
})
```

### Running the Server

To launch a basic server, create a `server.go` and populate it with the code of [server.go](examples/server/server.go). This example will guide you through setting up the server configuration, registering your custom function, and starting the server.

For more detailed examples, please refer to the `examples` folder:
- [Client Example](examples/client/client.go)
- [Server Example](examples/server/server.go)

## Contributing

We warmly welcome contributions from everyone. Your suggestions, code contributions, and feedback enhance the project for everyone. Engage with us and help make **Gofel** even better!

## License

**Gofel** is made available under the GNU GENERAL PUBLIC LICENSE. For more details, see the [LICENSE](LICENSE) file.

## Acknowledgments

- Thanks to the Gorilla WebSocket library.
- Gratitude to all contributors and the supportive community around the Gofel project.

We hope **Gofel** helps you build amazing applications. Happy coding!