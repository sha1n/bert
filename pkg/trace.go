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
	perceivedTime time.Duration
	usrCPUTime    time.Duration
	sysCPUTime    time.Duration
	error         error
}

func (t trace) ID() string {
	return t.id
}

func (t trace) PerceivedTime() time.Duration {
	return t.perceivedTime
}

func (t trace) SystemCPUTime() time.Duration {
	return t.sysCPUTime
}

func (t trace) UserCPUTime() time.Duration {
	return t.usrCPUTime
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

	return tr.endFn(t)
}

func (tr *tracer) endFn(t trace) api.End {
	return func(execInfo *api.ExecutionInfo, exitError error) {
		if execInfo != nil {
			t.perceivedTime, t.usrCPUTime, t.sysCPUTime = execInfo.PerceivedTime, execInfo.UserTime, execInfo.SystemTime
		}
		t.error = exitError

		tr.stream <- t
	}
}

func (tr *tracer) Stream() chan api.Trace {
	return tr.stream
}
