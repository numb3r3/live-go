package message

import (
	"github.com/numb3r3/h5-rtms-server/utils"
	"github.com/golang/snappy"
)

// Frame represents a message frame which is sent through the wire to the
// remote server and contains a set of messages.
type Frame []Message

// Message represents a message which has to be forwarded or stored.
type Message struct {
	Time    int64  `json:"ts,omitempty"`   // The timestamp of the message
	Ssid    Ssid   `json:"ssid,omitempty"` // The Ssid of the message
	Channel []byte `json:"chan,omitempty"` // The channel of the message
	Payload []byte `json:"data,omitempty"` // The payload of the message
	TTL     uint32 `json:"ttl,omitempty"`  // The time-to-live of the message
}

// Size returns the byte size of the message.
func (m *Message) Size() int64 {
	return int64(len(m.Payload))
}

// Encode encodes the message frame
func (f *Frame) Encode() (out []byte, err error) {
	// TODO: optimize
	var enc []byte
	if enc, err = utils.Encode(f); err == nil {
		out = snappy.Encode(out, enc)
		return
	}
	return
}

// Append appends the message to a frame.
func (f *Frame) Append(time int64, ssid Ssid, channel, payload []byte) {
	*f = append(*f, Message{Time: time, Ssid: ssid, Channel: channel, Payload: payload})
}

// DecodeFrame decodes the message frame from the decoder.
func DecodeFrame(buf []byte) (out Frame, err error) {
	// TODO: optimize
	var buffer []byte
	if buf, err = snappy.Decode(buffer, buf); err == nil {
		out = make(Frame, 0, 64)
		err = utils.Decode(buf, &out)
	}
	return
}