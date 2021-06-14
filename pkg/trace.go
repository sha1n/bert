package pkg

import (
	"time"

	"github.com/sha1n/bert/api"
)

type tracer struct {
	stream chan api.Trace
}

type trace struct {
	id            string
	cpuTimer      api.CPUTimer
	perceivedTime time.Duration
	usrTime       time.Duration
	sysTime       time.Duration
	error         error
}

func (t trace) ID() string {
	return t.id
}

func (t trace) Elapsed() time.Duration {
	return t.perceivedTime
}

func (t trace) System() time.Duration {
	return t.sysTime
}

func (t trace) User() time.Duration {
	return t.usrTime
}

func (t trace) Error() error {
	return t.error
}

func newTrace(id string) trace {
	return trace{
		id:       id,
		cpuTimer: NewChildrenCPUTimer(),
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
	t.cpuTimer.Start()

	return tr.endFn(t)
}

func (tr *tracer) endFn(t trace) api.End {
	return func(exitError error) {
		t.perceivedTime, t.usrTime, t.sysTime = t.cpuTimer.Elapsed()
		t.error = exitError

		tr.stream <- t
	}
}

func (tr *tracer) Stream() chan api.Trace {
	return tr.stream
}
