package broker

import (
	"net"
	"net/http"
	"time"

	// "github.com/gorilla/websocket"
	"github.com/numb3r3/h5-rtms-server/log"
	"github.com/numb3r3/h5-rtms-server/network/websocket"
	"github.com/numb3r3/h5-rtms-server/network/listener"
	"github.com/spf13/viper"
)

// Service represents the main structure.
type Service struct {
	Closing     chan bool    // The channel for closing signal.
	Config      *viper.Viper // The configuration for the service.
	http        *http.Server // The underlying HTTP server.
	startTime   time.Time    // The start time of the service.
	connections int64        // The number of currently open connections.
}

// NewService creates a new service.
func NewService(cfg *viper.Viper) (s *Service, err error) {
	s = &Service{
		Closing: make(chan bool),
		Config:  cfg,
		http:    new(http.Server),
	}

	// Create a new HTTP request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.onRequest)

	// Attach handlers
	s.http.Handler = mux

	return s, nil
}

// Listen starts the service.
func (s *Service) Listen() (err error) {
	defer s.Close()

	// Setup the listeners on both default and a secure addresses
	s.listen(s.Config.GetString("listen_addr"))

	// Set the start time and report status
	s.startTime = time.Now().UTC()
	logging.Info("service", "service started")

	// Block
	select {}
}

// listen configures an main listener on a specified address.
func (s *Service) listen(address string) {
	logging.Info("service", "starting the listener", address)

	l, err := listener.NewListener(address)
	if err != nil {
		panic(err)
	}

	// // Set the read timeout on our mux listener
	// l.SetReadTimeout(120 * time.Second)

	// // Configure the matchers
	// l.ServeAsync(s.http.Serve)
	// // l.ServeAsync(listener.MatchAny(), s.tcp.Serve)
	// go l.Serve()
}

// Occurs when a new client connection is accepted.
func (s *Service) onAcceptConn(t net.Conn) {
	conn := s.newConn(t)
	go conn.Process()
}

// Occurs when a new HTTP request is received.
func (s *Service) onRequest(w http.ResponseWriter, r *http.Request) {
	if ws, ok := websocket.TryUpgrade(w, r); ok {
		s.onAcceptConn(ws)
		return
	}
}

// Close closes gracefully the service.,
func (s *Service) Close() {

	// Notify we're closed
	close(s.Closing)
}
