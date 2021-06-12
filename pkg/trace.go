package pkg

import (
	"time"

	"github.com/sha1n/benchy/api"
)

type tracer struct {
	stream chan api.Trace
}

type trace struct {
	id            string
	cpuTimer      CPUTimer
	start         time.Time
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
		cpuTimer: newChildrenCPUTimer(),
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
	t.start = time.Now()

	return tr.endFn(t)
}

func (tr *tracer) endFn(t trace) api.End {
	return func(exitError error) {
		t.usrTime, t.sysTime = t.cpuTimer.Elapsed()
		t.perceivedTime = time.Since(t.start)
		t.error = exitError

		tr.stream <- t
	}
}

func (tr *tracer) Stream() chan api.Trace {
	return tr.stream
}

// CPUTimer an abstraction for a platform system CPU timer
type CPUTimer interface {
	Start() func() (time.Duration, time.Duration)
	Elapsed() (usr time.Duration, sys time.Duration)
}

// func newChildrenCPUTimer() CPUTimer {
// 	switch runtime.GOOS {
// 	case "windows":
// 		return NoopCPUTimer{}
// 	default:

// 		return newUnixChildrenCPUTimer()
// 	}
// }

// // NoopCPUTimer NOOP implementation of CPUTimer
// type NoopCPUTimer struct{}

// // Start ...
// func (t NoopCPUTimer) Start() func() (time.Duration, time.Duration) {
// 	return t.Elapsed
// }

// // Elapsed return 0, 0
// func (t NoopCPUTimer) Elapsed() (usr time.Duration, sys time.Duration) {
// 	return time.Nanosecond * 0, time.Nanosecond * 0
// }
