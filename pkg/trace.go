package pkg

import (
	"syscall"
	"time"

	"github.com/sha1n/benchy/api"
)

type tracer struct {
	stream chan api.Trace
}

type trace struct {
	id            string
	cpuTimer      *childrenCPUTimer
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

type childrenCPUTimer struct {
	r            syscall.Rusage
	sysTimeStart time.Time
	usrTimeStart time.Time
}

func newChildrenCPUTimer() *childrenCPUTimer {
	return &childrenCPUTimer{}
}

func (t *childrenCPUTimer) Start() func() (time.Duration, time.Duration) {
	err := syscall.Getrusage(syscall.RUSAGE_CHILDREN, &t.r)
	if err != nil {
		panic(err)
	}

	t.sysTimeStart = time.Unix(t.r.Stime.Unix())
	t.usrTimeStart = time.Unix(t.r.Utime.Unix())

	return t.Elapsed
}

func (t *childrenCPUTimer) Elapsed() (usr time.Duration, sys time.Duration) {
	err := syscall.Getrusage(syscall.RUSAGE_CHILDREN, &t.r)
	if err != nil {
		panic(err)
	}

	sysTimeEnd := time.Unix(t.r.Stime.Unix())
	usrTimeEnd := time.Unix(t.r.Utime.Unix())

	return usrTimeEnd.Sub(t.usrTimeStart), sysTimeEnd.Sub(t.sysTimeStart)
}
