package broker

import (
	"bufio"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/numb3r3/h5-rtms-server/log"
)

// Conn represents an incoming connection.
type Conn struct {
	sync.Mutex
	tracked  uint32
	socket   net.Conn
	username string
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
	reader := bufio.NewReaderSize(c.socket, 65536)

	for {
		// Set read/write deadlines so we can close dangling connections
		c.socket.SetDeadline(time.Now().Add(time.Second * 120))

		// Decode an incoming package

	}
}

// Close terminates the connection.
func (c *Conn) Close() error {
	logging.Info("connection closed.")

	// Unsubscribe from everything, no need to lock since each Unsubscribe is
	// already locked. Locking the 'Close()' would result in a deadlock.
	for _, counter := range c.subs.All() {
		c.service.onUnsubscribe(counter.Ssid, c)
		c.service.notifyUnsubscribe(c, counter.Ssid, counter.Channel)
	}

	// Attempt to recover a panic
	if r := recover(); r != nil {
		logging.Info("closing", fmt.Sprintf("pancic recovered: %s \n %s", r, debug.Stack()))
	}

	// Close the transport and decrement the connection counter
	atomic.AddInt64(&c.service.connections, -1)
	return c.socket.Close()
}
