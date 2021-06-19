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
	startTime     time.Time
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
		id: id,
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
	t.startTime = time.Now()

	return tr.endFn(t)
}

func (tr *tracer) endFn(t trace) api.End {
	return func(execInfo *api.ExecutionInfo, exitError error) {
		t.perceivedTime = time.Since(t.startTime)
		t.usrTime, t.sysTime = execInfo.UserTime, execInfo.SystemTime
		t.error = exitError

		tr.stream <- t
	}
}

func (tr *tracer) Stream() chan api.Trace {
	return tr.stream
}
