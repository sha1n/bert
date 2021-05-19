package pkg

import (
	"github.com/sha1n/benchy/api"
	"time"
)

type tracer struct {
	stream chan api.Trace
}

type trace struct {
	id      string
	start   time.Time
	elapsed time.Duration
	error   error
}

func (t *trace) end(exitError error) {
	t.elapsed = time.Since(t.start)
	t.error = exitError
}

func (t *trace) ID() string {
	return t.id
}

func (t *trace) Elapsed() time.Duration {
	return t.elapsed
}

func (t *trace) Error() error {
	return t.error
}

func newTrace(id string) *trace {
	return &trace{
		id:    id,
		start: time.Now(),
	}
}

// NewTracer creates a new tracer
func NewTracer(bufferSize int) api.Tracer {
	return &tracer{
		stream: make(chan api.Trace, bufferSize),
	}
}

func (tr *tracer) Start(i api.Identifiable) api.End {
	t := newTrace(i.ID())

	return tr.endFn(t)
}

func (tr *tracer) endFn(t *trace) api.End {
	return func(exitError error) {
		t.elapsed = time.Since(t.start)
		t.error = exitError

		tr.stream <- t
	}
}

func (tr *tracer) Stream() chan api.Trace {
	return tr.stream
}
