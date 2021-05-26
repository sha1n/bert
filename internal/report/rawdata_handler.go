package report

import "github.com/sha1n/benchy/api"

// RawDataHandler an abstraction for raw data report trace event handlers.
type RawDataHandler interface {
	// Handle receives a trace event and processes it immediately
	Handle(api.Trace) error
}
