package broker

import (
	// "io"
	// "bufio"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/numb3r3/live-go/log"
)

// Conn represents an incoming connection.
type Conn struct {
	sync.Mutex
	tracked  uint32
	socket   net.Conn
	username string
	service  *Service // The service for this connection.
	guid     string
}

// NewConn creates a new connection.
func (s *Service) newConn(t net.Conn) *Conn {
	c := &Conn{
		tracked: 0,
		service: s,
		socket:  t,
	}

	// TODO: generate a global unique id

	logging.Info("net connection created.")

	// Increment the connection counter
	atomic.AddInt64(&s.connections, 1)
	return c
}

// Process processes the messages.
func (c *Conn) Process() error {
	defer c.Close()
	// reader := bufio.NewReaderSize(c.socket, 65536)

	for {
		// Set read/write deadlines so we can close dangling connections
		c.socket.SetDeadline(time.Now().Add(time.Second * 120))

		// Decode an incoming package

		// b := make([]byte, 1)
		// if _, err := io.ReadFull(reader, b); err != nil {
		// 	return nil, 0, 0, err
		// }

	}
}

// Close terminates the connection.
func (c *Conn) Close() error {
	logging.Info("connection closed.")

	// Attempt to recover a panic
	if r := recover(); r != nil {
		logging.Info("closing", fmt.Sprintf("pancic recovered: %s \n %s", r, debug.Stack()))
	}

	// Close the transport and decrement the connection counter
	atomic.AddInt64(&c.service.connections, -1)
	return c.socket.Close()
}
