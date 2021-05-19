package pkg

import (
	"sync"

	"github.com/sha1n/benchy/api"
)

// TraceSink is a Trace accumulator capable of subscribing to a TraceStream and provide access to accumulated events
type TraceSink struct {
	traces     map[api.ID][]api.Trace
	subscriber *StreamSubscriber
	mx         *sync.RWMutex
}

// NewTraceSink creates a new TraceSink for the specified stream.
func NewTraceSink(stream api.TraceStream) *TraceSink {
	s := &TraceSink{
		traces: make(map[api.ID][]api.Trace),
		mx:     &sync.RWMutex{},
	}

	s.subscriber = NewStreamSubscriber(stream, s.add)
	return s
}

// Subscribe subscribes to the trace events stream and returns an unsubscribe handler.
// Following that call, events start to accumulate.
func (s *TraceSink) Subscribe() Unsubscribe {
	return s.subscriber.Subscribe()
}

// Summary returns an Summary object containing stats about accumulated events so far.
func (s *TraceSink) Summary() api.Summary {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return NewSummary(s.traces)
}

func (s *TraceSink) add(trace api.Trace) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.traces[trace.ID()] == nil {
		s.traces[trace.ID()] = []api.Trace{}
	}

	s.traces[trace.ID()] = append(s.traces[trace.ID()], trace)

	return nil
}
