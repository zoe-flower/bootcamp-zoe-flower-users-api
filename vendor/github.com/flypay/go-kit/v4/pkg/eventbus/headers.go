package eventbus

import "time"

const (
	// HeaderTimestamp contains the time at which the message was emitted in time.RFC3339Nano
	HeaderTimestamp = "timestamp"
	// HeaderAttempt contains the attempt number, the first time an event is seen the value will be `1`.
	HeaderAttempt = "attempt"
)

// Header is a metadata value that should be sent with the event
type Header struct {
	Key   string
	Value string
}

// EmittedAt returns the time at which an event was emitted onto the eventbus.
func EmittedAt(headers ...Header) (time.Time, bool) {
	var timestamp string
	for _, h := range headers {
		if h.Key == HeaderTimestamp {
			timestamp = h.Value
			break
		}
	}

	if timestamp == "" {
		return time.Time{}, false
	}

	happenedAt, err := time.Parse(time.RFC3339Nano, timestamp)
	return happenedAt, err == nil
}

// EmittedUnixNano returns the unix timestamp, with nano precision, at which an event was emitted onto the eventbus
func EmittedUnixNano(headers ...Header) (int64, bool) {
	at, ok := EmittedAt(headers...)
	return at.UnixNano(), ok
}

// WithEmittedAt adds the current timestamp header for later use by `HappenedAt`.
func WithEmittedAt(at time.Time, headers ...Header) []Header {
	return append(headers, Header{
		Key:   HeaderTimestamp,
		Value: at.Format(time.RFC3339Nano),
	})
}
