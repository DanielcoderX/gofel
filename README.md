# Gofel

**Gofel** is a high-performance, lightweight RPC (Remote Procedure Call) framework crafted in Go, aimed at streamlining the development of real-time communication servers. By harnessing the capabilities of the Gorilla WebSocket library, **Gofel** empowers developers to create scalable, efficient WebSocket applications with ease.

## Features

- **Easy Setup**: Quickly launch your server with minimal configuration.
- **Customizable Configurations**: Adapt the server settings to perfectly fit your application requirements.
- **WebSocket Integration**: Built on the reliable and high-performance Gorilla WebSocket library.
- **Concurrent Processing**: Designed for handling multiple requests simultaneously, enhancing real-time communication efficiency.
- **Lightweight Framework**: Focus more on your application's features and less on managing boilerplate code.
- **MessagePack and JSON Support**: Support for both MessagePack and JSON data formats, allowing for efficient data transmission.
## Getting Started

Embark on your Gofel journey with these simple steps:

### Prerequisites

Make sure you have Go installed (version 1.15 or newer). Download it from [Go Downloads](https://golang.org/dl/).

### Installation

Incorporate Gofel into your project by running the following command:

```go
go get github.com/DanielcoderX/gofel
```

### Function Registration and Message Handling

When setting up your server with Gofel, you'll typically register functions that can handle incoming requests in the following format:

```json
{
    "function": "echo",
    "data": "userdatatofunction"
}
```

Each registered function must be capable of parsing the following structure to perform the required action. Here's an example of how to register such a function:

```go
server.On("echo", func(conn *wsconn.ConnectionWrapper, data interface{}) {
    err := api.SendResponse(conn, data)
    if err != nil {
        log.Printf("Failed to send response: %v", err)
    }
})
```

### Running the Server or Client

To launch a basic server or Client, refer to the [examples](examples) directory for a step-by-step guide on how to set up and use the server or client.

In the examples directory, you'll find a working setup that you can use as a reference to get started.

## Contributing

We warmly welcome contributions from everyone. Your suggestions, code contributions, and feedback enhance the project for everyone. Engage with us and help make **Gofel** even better!

## License

**Gofel** is made available under the GNU GENERAL PUBLIC LICENSE. For more details, see the [LICENSE](LICENSE) file.

## Acknowledgments

- Thanks to the Gorilla WebSocket library.
