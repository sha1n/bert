package pkg

import (
	"context"
	"sync"

	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// Unsubscribe drains a TraceStream from buffered events and unsubscribes from it.
type Unsubscribe = func()

// HandleFn handles a trace event
type HandleFn = func(api.Trace) error

// StreamSubscriber FIXME
type StreamSubscriber struct {
	stream   chan api.Trace
	handleFn HandleFn
}

// NewStreamSubscriber FIXME
func NewStreamSubscriber(stream api.TraceStream, handleFn HandleFn) *StreamSubscriber {
	return &StreamSubscriber{
		stream:   stream,
		handleFn: handleFn,
	}
}

// Subscribe FIXME
func (s *StreamSubscriber) Subscribe() Unsubscribe {
	context, cancel := context.WithCancel(context.Background())
	startWaitGroup := &sync.WaitGroup{}
	startWaitGroup.Add(1)

	drainWaitGroup := &sync.WaitGroup{}
	drainWaitGroup.Add(1)

	go func() {
		startWaitGroup.Done()

		for {
			select {
			case <-context.Done():
				s.drain()
				drainWaitGroup.Done()

				return

			case trace := <-s.stream:
				s.handle(trace)
			}
		}
	}()

	startWaitGroup.Wait()

	return func() {
		cancel()
		drainWaitGroup.Wait()
	}
}

func (s *StreamSubscriber) drain() {
	for len(s.stream) > 0 {
		s.handle(<-s.stream)
	}
}

func (s *StreamSubscriber) handle(trace api.Trace) {
	if err := s.handleFn(trace); err != nil {
		log.Errorf("Failed to handle trace event. Error: %v", err)
	}
}
