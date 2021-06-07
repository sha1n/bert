package pkg

import (
	"context"
	"fmt"
	"reflect"
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
	stream         chan api.Trace
	handleFn       HandleFn
	mx             *sync.RWMutex
	drainWaitGroup *sync.WaitGroup
}

// NewStreamSubscriber returns a new subscriber for the specified TraceStream. Does not subscribe immediately.
func NewStreamSubscriber(stream api.TraceStream, handleFn HandleFn) *StreamSubscriber {
	return &StreamSubscriber{
		stream:   stream,
		handleFn: handleFn,
		mx:       &sync.RWMutex{},
	}
}

// Subscribe subscribes to the TraceStream and starts handling events.
// Returns an Unsubscribe handler that drains the underlaying channel and stops consuming events.
// This method can only be called once before Unsubscribe is called.
func (s *StreamSubscriber) Subscribe() Unsubscribe {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.drainWaitGroup != nil {
		panic(fmt.Errorf("programmer error: %s can only subscribe one", reflect.TypeOf(s)))
	}

	context, cancel := context.WithCancel(context.Background())
	startWaitGroup := &sync.WaitGroup{}
	startWaitGroup.Add(1)

	s.drainWaitGroup = &sync.WaitGroup{}
	s.drainWaitGroup.Add(1)

	go func() {
		startWaitGroup.Done()

		for {
			select {
			case <-context.Done():
				s.drain()
				s.drainWaitGroup.Done()

				return

			case trace := <-s.stream:
				s.handle(trace)
			}
		}
	}()

	// we don't return before the subscription routine starts.
	startWaitGroup.Wait()

	return func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		cancel()
		s.drainWaitGroup.Wait()
		s.drainWaitGroup = nil
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
