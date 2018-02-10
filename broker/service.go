package broker

import (

	"github.com/spf13/viper"
	"github.com/numb3r3/h5-rtms-server/log"
)

// Service represents the main structure.
type Service struct {
	Closing       chan bool                 // The channel for closing signal.
	Config        *viper.Viper            	// The configuration for the service.
	subscriptions *message.Trie             // The subscription matching trie.
	http          *http.Server              // The underlying HTTP server.
	tcp           *tcp.Server               // The underlying TCP server.
	cluster       *cluster.Swarm            // The gossip-based cluster mechanism.
	startTime     time.Time                 // The start time of the service.
	presence      chan *presenceNotify      // The channel for presence notifications.
	// querier       *QueryManager             // The generic query manager.
	// contracts     security.ContractProvider // The contract provider for the service.
	// storage       storage.Storage           // The storage provider for the service.
	// metering      usage.Metering            // The usage storage for metering contracts.
	connections   int64                     // The number of currently open connections.
}

// NewService creates a new service.
func NewService(cfg *viper.Viper) (s *Service, err error) {
	s = &Service{
		Closing:       make(chan bool),
		Config:        cfg,
		// subscriptions: message.NewTrie(),
		http:          new(http.Server),
		tcp:           new(tcp.Server),
		// presence:      make(chan *presenceNotify, 100),
		// storage:       new(storage.Noop),
	}

	// Create a new HTTP request multiplexer
	mux := http.NewServeMux()
	// mux.HandleFunc("/health", s.onHealth)
	// mux.HandleFunc("/keygen", s.onHTTPKeyGen)
	// mux.HandleFunc("/presence", s.onHTTPPresence)
	// mux.HandleFunc("/debug/pprof/", pprof.Index)          // TODO: use config flag to enable/disable this
	// mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline) // TODO: use config flag to enable/disable this
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile) // TODO: use config flag to enable/disable this
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)   // TODO: use config flag to enable/disable this
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)     // TODO: use config flag to enable/disable this
	mux.HandleFunc("/", s.onRequest)

	// Attach handlers
	s.http.Handler = mux
	s.tcp.OnAccept = s.onAcceptConn
	
	// s.querier = newQueryManager(s)
	return s, nil
}
