package exec

import (
	"math/rand"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/gommons/pkg/test"

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
	stream := make(chan api.Trace)
	subscriber := NewStreamSubscriber(stream, func(t api.Trace) error { return nil })

	unsubscribe := subscriber.Subscribe()
	defer unsubscribe()

	assert.Panics(t, func() {
		subscriber.Subscribe()
	})
}

func createTracerWithNBufferedEvents(n int) api.Tracer {
	tracer := NewTracer(n)
	spec := aSpec(true, 1)

	for i := 0; i < n; i++ {
		tracer.Start(spec.Scenarios[0])(&api.ExecutionInfo{}, nil)
	}

	return tracer
}

func aSpec(alternate bool, executions int) api.BenchmarkSpec {
	return api.BenchmarkSpec{
		Executions: executions,
		Alternate:  alternate,
		Scenarios: []api.ScenarioSpec{
			{
				Name: test.RandomString(),
				Command: &api.CommandSpec{
					Cmd: []string{"cmd"},
				},
			},
		},
	}
}
