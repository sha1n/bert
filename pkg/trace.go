package pkg

import "time"

type End = func(error)
type Id = string

type Identifiable interface {
	Id() Id
}

type Tracer interface {
	Start(i Identifiable) End
	Summary() TracerSummary
}

type Trace interface {
	Id() string
	Elapsed() time.Duration
	Error() error
}

type tracer struct {
	traces map[Id][]Trace
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

func (t *trace) Id() string {
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

func NewTracer() Tracer {
	return &tracer{
		traces: make(map[Id][]Trace),
	}
}

func (tr *tracer) Start(i Identifiable) End {
	t := newTrace(i.Id())

	if tr.traces[i.Id()] == nil {
		tr.traces[i.Id()] = []Trace{}
	}

	tr.traces[i.Id()] = append(tr.traces[i.Id()], t)

	return t.end
}

func (tr *tracer) Summary() TracerSummary {
	return NewTracerSummary(tr.traces)
}
