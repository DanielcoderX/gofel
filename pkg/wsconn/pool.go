package wsconn

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionPool manages a pool of WebSocket connections.
type ConnectionPool struct {
	pool  []*ConnectionWrapper // Array of connections
	mutex sync.Mutex           // Mutex for thread safety
}

var (
	GlobalPool *ConnectionPool // Global connection pool instance
)

// InitConnectionPool initializes a new ConnectionPool with a given capacity.
func InitConnectionPool(capacity int) {
	GlobalPool = &ConnectionPool{
		pool: make([]*ConnectionWrapper, capacity),
	}
	for i := 0; i < capacity; i++ {
		GlobalPool.pool[i] = &ConnectionWrapper{}
	}
}

// GetConnection retrieves a websocket connection from the pool.
func (cp *ConnectionPool) GetConnection(conn *websocket.Conn) (*ConnectionWrapper, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for _, c := range cp.pool {
		if c == nil {
			continue
		}
		if c.IsAvailable() {
			c.Reset(conn)
			c.SetInUse(true) // Mark as in use
			return c, nil
		}
	}
	return nil, errors.New("no available connections in the pool")
}

// ReleaseConnection releases a websocket connection back to the pool.
func (cp *ConnectionPool) ReleaseConnection(conn *ConnectionWrapper) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for _, c := range cp.pool {
		if c == conn {
			c.SetInUse(false) // Mark as available
			return
		}
	}
}

// AddConnection adds a new websocket connection to the pool.
func (cp *ConnectionPool) AddConnection(conn *ConnectionWrapper) error {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for i, c := range cp.pool {
		if c == nil {
			cp.pool[i] = conn
			return nil
		}
	}
	return errors.New("connection pool is full")
}

// RemoveConnection removes a websocket connection from the pool.
func (cp *ConnectionPool) RemoveConnection(conn *ConnectionWrapper) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for i, c := range cp.pool {
		if c == conn {
			cp.pool[i] = nil
			return
		}
	}
}
