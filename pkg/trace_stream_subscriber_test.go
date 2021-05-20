package pkg

import (
	"math/rand"
	"testing"

	"github.com/sha1n/benchy/api"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	eventCount := rand.Intn(100) + 100 // using a large number to make sure all buffered events are handled when unsubscribing
	tracer := createTracerWithNBufferedEvents(eventCount)
	received := []api.Trace{}
	handleFn := func(t api.Trace) error {
		received = append(received, t)
		return nil
	}

	subscriber := NewStreamSubscriber(tracer.Stream(), handleFn)
	unsubscribe := subscriber.Subscribe()

	unsubscribe()
	assert.Equal(t, eventCount, len(received))
}

func TestCanOnlySubscribeOnce(t *testing.T) {
	expectPanic := func() {
		v := recover()
		assert.NotNil(t, v)
	}
	defer expectPanic()

	stream := make(chan api.Trace)
	subscriber := NewStreamSubscriber(stream, func(t api.Trace) error { return nil })

	unsubscribe := subscriber.Subscribe()
	defer unsubscribe()

	subscriber.Subscribe() // => expecting panic
}

func createTracerWithNBufferedEvents(n int) api.Tracer {
	tracer := NewTracer(n)
	spec := exampleSpec()

	for i := 0; i < n; i++ {
		tracer.Start(spec.Scenarios[0])(nil)
	}

	return tracer
}