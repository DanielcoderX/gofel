package wsconn

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionWrapper wraps a WebSocket connection with a mutex for synchronized writes.
type ConnectionWrapper struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// NewConnectionWrapper creates a new ConnectionWrapper.
func NewConnectionWrapper(conn *websocket.Conn) *ConnectionWrapper {
	return &ConnectionWrapper{conn: conn}
}

// SendMessage safely sends a message through the WebSocket connection.
func (cw *ConnectionWrapper) SendMessage(messageType int, data []byte) error {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.conn.WriteMessage(messageType, data)
}

// ReadMessage reads a message from the WebSocket connection.
func (cw *ConnectionWrapper) ReadMessage() (messageType int, p []byte, err error) {
	return cw.conn.ReadMessage()
}

// Close closes the WebSocket connection.
func (cw *ConnectionWrapper) Close() error {
	return cw.conn.Close()
}
