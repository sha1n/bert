package pkg

import "time"

// End ends a trace
type End = func(error)

// ID well..
type ID = string

// Identifiable an abstraction for identifiable objects
type Identifiable interface {
	ID() ID
}

// Tracer a global tracing handler that accumulates trace data and provides access to it.
type Tracer interface {
	Start(i Identifiable) End
	Summary() TracerSummary
}

// Trace a single time trace
type Trace interface {
	ID() string
	Elapsed() time.Duration
	Error() error
}

type tracer struct {
	traces map[ID][]Trace
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
func NewTracer() Tracer {
	return &tracer{
		traces: make(map[ID][]Trace),
	}
}

func (tr *tracer) Start(i Identifiable) End {
	t := newTrace(i.ID())

	if tr.traces[i.ID()] == nil {
		tr.traces[i.ID()] = []Trace{}
	}

	tr.traces[i.ID()] = append(tr.traces[i.ID()], t)

	return t.end
}

func (tr *tracer) Summary() TracerSummary {
	return NewTracerSummary(tr.traces)
}
