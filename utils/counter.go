package utils

// Counters represents a subscription counting map.
type Counters struct {
	sync.Mutex
	m map[uint32]*Counter
}

// Counter represents a single subscription counter.
type Counter struct {
	Ssid    Ssid
	Channel []byte
	Counter int
}

// NewCounters creates a new container.
func NewCounters() *Counters {
	return &Counters{
		m: make(map[uint32]*Counter),
	}
}

// Increment increments the subscription counter.
func (s *Counters) Increment(ssid Ssid, channel []byte) (first bool) {
	s.Lock()
	defer s.Unlock()

	m := s.getOrCreate(ssid, channel)
	m.Counter++
	return m.Counter == 1
}

// Decrement decrements a subscription counter.
func (s *Counters) Decrement(ssid Ssid) (last bool) {
	s.Lock()
	defer s.Unlock()

	key := ssid.GetHashCode()
	if m, exists := s.m[key]; exists {
		m.Counter--

		// Remove if there's no subscribers left
		if m.Counter <= 0 {
			delete(s.m, ssid.GetHashCode())
			return true
		}
	}

	return false
}

// All returns all counters.
func (s *Counters) All() []Counter {
	s.Lock()
	defer s.Unlock()

	clone := make([]Counter, 0, len(s.m))
	for _, m := range s.m {
		clone = append(clone, *m)
	}

	return clone
}

// getOrCreate retrieves a single subscription meter or creates a new one.
func (s *Counters) getOrCreate(ssid Ssid, channel []byte) (meter *Counter) {
	key := ssid.GetHashCode()
	if m, exists := s.m[key]; exists {
		return m
	}

	meter = &Counter{
		Ssid:    ssid,
		Channel: channel,
		Counter: 0,
	}
	s.m[key] = meter
	return
}