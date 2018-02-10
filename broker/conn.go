package broker

import (
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/numb3r3/h5-rtms-server/broker/message"
	"github.com/numb3r3/h5-rtms-server/log"
)


// Conn represents an incoming connection.
type Conn struct {
	sync.Mutex
	tracked  uint32            // Whether the connection was already tracked or not.
	socket   net.Conn          // The transport used to read and write messages.
	// luid     security.ID       // The locally unique id of the connection.
	guid     string            // The globally unique id of the connection.
	service  *Service          // The service for this connection.
	subs     *message.Counters // The subscriptions for this connection.
}