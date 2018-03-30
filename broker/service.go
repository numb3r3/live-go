package broker

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/numb3r3/live-go/log"
	"github.com/numb3r3/live-go/network/listener"
	"github.com/numb3r3/live-go/network/websocket"
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
	mux.HandleFunc("/health", s.onHealth)
	mux.HandleFunc("/", s.onRequest)

	// Attach handlers
	s.http.Handler = mux

	return s, nil
}

// Listen starts the service.
func (s *Service) Listen() (err error) {
	defer s.Close()
	s.hookSignals()

	// Setup the listeners on both default and a secure addresses
	s.listen(s.Config.GetString("listen_addr"))

	// Set the start time and report status
	s.startTime = time.Now().UTC()
	logging.Info("service started")

	// Block
	select {}
}

// listen configures an main listener on a specified address.
func (s *Service) listen(address string) {
	logging.Info("starting the listener", address)

	l, err := listener.NewListener(address)
	if err != nil {
		panic(err)
	}

	// Set the read timeout on our mux listener
	l.SetReadTimeout(120 * time.Second)

	l.ServeAsync(s.http.Serve)

	// l.ServeAsync(listener.MatchAny(), s.tcp.Serve)
	go l.Serve()
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

// Occurs when a new HTTP health check is received.
func (s *Service) onHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

// OnSignal will be called when a OS-level signal is received.
func (s *Service) onSignal(sig os.Signal) {
	switch sig {
	case syscall.SIGTERM:
		fallthrough
	case syscall.SIGINT:
		logging.Infof("received signal %s, exiting...", sig.String())
		s.Close()
		os.Exit(0)
	}
}

// OnSignal starts the signal processing and makes su
func (s *Service) hookSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			s.onSignal(sig)
		}
	}()
}

// Close closes gracefully the service.,
func (s *Service) Close() {

	// Notify we're closed
	close(s.Closing)
}
