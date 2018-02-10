package listener

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type processor struct {
	// matchers []Matcher
	listen   muxListener
}

// Listener represents a listener used for multiplexing protocols.
type Listener struct {
	root         net.Listener
	bufferSize   int
	errorHandler ErrorHandler
	closing      chan struct{}
	processor    muxListener
	readTimeout  time.Duration
}

// Accept waits for and returns the next connection to the listener.
func (m *Listener) Accept() (net.Conn, error) {
	return m.root.Accept()
}

// ServeAsync adds a protocol based on the matcher and serves it.
func (m *Listener) ServeAsync(serve func(l net.Listener) error) {
	ml := muxListener{
		Listener:    m.root,
		connections: make(chan net.Conn, m.bufferSize),
	}
	m.processor = ml
	go serve(ml)
}

// Serve starts multiplexing the listener.
func (m *Listener) Serve() error {
	var wg sync.WaitGroup

	defer func() {
		close(m.closing)
		wg.Wait()

		close(processor.connections)
		for c := range processor.connections {
			_ = c.Close()
		}
	}()

	for {
		c, err := m.root.Accept()
		if err != nil {
			if !m.handleErr(err) {
				return err
			}
			continue
		}

		wg.Add(1)
		go m.serve(c, m.closing, &wg)
	}
}

func (m *Listener) serve(c net.Conn, donec <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	muc := newConn(c)
	if m.readTimeout > noTimeout {
		_ = c.SetReadDeadline(time.Now().Add(m.readTimeout))
	}
	for _, sl := range m.matchers {
		for _, processor := range sl.matchers {
			matched := processor(muc.startSniffing())
			if matched {
				muc.doneSniffing()
				if m.readTimeout > noTimeout {
					_ = c.SetReadDeadline(time.Time{})
				}
				select {
				case sl.listen.connections <- muc:
				case <-donec:
					_ = c.Close()
				}
				return
			}
		}
	}

	_ = c.Close()
	err := ErrNotMatched{c: c}
	if !m.handleErr(err) {
		_ = m.root.Close()
	}
}

// HandleError registers an error handler that handles listener errors.
func (m *Listener) HandleError(h ErrorHandler) {
	m.errorHandler = h
}

func (m *Listener) handleErr(err error) bool {
	if !m.errorHandler(err) {
		return false
	}

	if ne, ok := err.(net.Error); ok {
		return ne.Temporary()
	}

	return false
}

// Close closes the listener
func (m *Listener) Close() error {
	return m.root.Close()
}

// Addr returns the listener's network address.
func (m *Listener) Addr() net.Addr {
	return m.root.Addr()
}


// ------------------------------------------------------------------------------------

type muxListener struct {
	net.Listener
	connections chan net.Conn
}

func (l muxListener) Accept() (net.Conn, error) {
	c, ok := <-l.connections
	if !ok {
		return nil, ErrListenerClosed
	}
	return c, nil
}


// ------------------------------------------------------------------------------------

// Conn wraps a net.Conn and provides transparent sniffing of connection data.
type Conn struct {
	net.Conn
	buffer sniffer
}

// NewConn creates a new sniffed connection.
func newConn(c net.Conn) *Conn {
	return &Conn{
		Conn:   c,
		buffer: sniffer{source: c},
	}
}

// Read reads the block of data from the underlying buffer.
func (m *Conn) Read(p []byte) (int, error) {
	return m.buffer.Read(p)
}

func (m *Conn) startSniffing() io.Reader {
	m.buffer.reset(true)
	return &m.buffer
}

func (m *Conn) doneSniffing() {
	m.buffer.reset(false)
}



// ------------------------------------------------------------------------------------

// Sniffer represents a io.Reader which can peek incoming bytes and reset back to normal.
type sniffer struct {
	source     io.Reader
	buffer     bytes.Buffer
	bufferRead int
	bufferSize int
	sniffing   bool
	lastErr    error
}

// Read reads data from the buffer.
func (s *sniffer) Read(p []byte) (int, error) {
	if s.bufferSize > s.bufferRead {
		bn := copy(p, s.buffer.Bytes()[s.bufferRead:s.bufferSize])
		s.bufferRead += bn
		return bn, s.lastErr
	} else if !s.sniffing && s.buffer.Cap() != 0 {
		s.buffer = bytes.Buffer{}
	}

	sn, sErr := s.source.Read(p)
	if sn > 0 && s.sniffing {
		s.lastErr = sErr
		if wn, wErr := s.buffer.Write(p[:sn]); wErr != nil {
			return wn, wErr
		}
	}
	return sn, sErr
}

// Reset resets the buffer.
func (s *sniffer) reset(snif bool) {
	s.sniffing = snif
	s.bufferRead = 0
	s.bufferSize = s.buffer.Len()
}