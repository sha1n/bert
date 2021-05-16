package pkg

import (
	"github.com/sha1n/benchy/api"
	"time"
)

type tracer struct {
	traces map[api.ID][]api.Trace
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
func NewTracer() api.Tracer {
	return &tracer{
		traces: make(map[api.ID][]api.Trace),
	}
}

func (tr *tracer) Start(i api.Identifiable) api.End {
	t := newTrace(i.ID())

	if tr.traces[i.ID()] == nil {
		tr.traces[i.ID()] = []api.Trace{}
	}

	tr.traces[i.ID()] = append(tr.traces[i.ID()], t)

	return t.end
}

func (tr *tracer) Summary() api.Summary {
	return NewSummary(tr.traces)
}
