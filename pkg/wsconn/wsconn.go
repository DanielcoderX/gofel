package wsconn

import (
    "github.com/gorilla/websocket"
    "sync"
)

// ConnectionWrapper wraps a WebSocket connection with a mutex for synchronized writes.
type ConnectionWrapper struct {
    conn  *websocket.Conn
    mu    sync.Mutex
    inUse bool // Indicates if the connection is currently in use
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

// Reset resets the ConnectionWrapper with a new connection.
func (cw *ConnectionWrapper) Reset(conn *websocket.Conn) {
    cw.mu.Lock()
    defer cw.mu.Unlock()
    cw.conn = conn
}

// IsAvailable checks if the connection is available in the pool.
func (cw *ConnectionWrapper) IsAvailable() bool {
    cw.mu.Lock()
    defer cw.mu.Unlock()
    return !cw.inUse
}

// SetInUse sets the inUse flag.
func (cw *ConnectionWrapper) SetInUse(inUse bool) {
    cw.mu.Lock()
    defer cw.mu.Unlock()
    cw.inUse = inUse
}
