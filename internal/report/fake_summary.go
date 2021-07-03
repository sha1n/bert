package report

import (
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg/exec"
)

type fakeTrace struct {
	id            string
	perceivedTime time.Duration
	usrCPUTime    time.Duration
	sysCPUTime    time.Duration
	error         error
}

func (t fakeTrace) ID() string {
	return t.id
}

func (t fakeTrace) PerceivedTime() time.Duration {
	return t.perceivedTime
}

func (t fakeTrace) SystemCPUTime() time.Duration {
	return t.sysCPUTime
}

func (t fakeTrace) UserCPUTime() time.Duration {
	return t.usrCPUTime
}

func (t fakeTrace) Error() error {
	return t.error
}

// NewFakeTrace creates a fake trace with the specified data
func NewFakeTrace(id string, elapsed, userTime, sysTime time.Duration, err error) api.Trace {
	return &fakeTrace{
		id:            id,
		perceivedTime: elapsed,
		usrCPUTime:    userTime,
		sysCPUTime:    sysTime,
		error:         err,
	}
}

// NewFakeSummary creates a new fake summary object with the specified trace events
func NewFakeSummary(traces ...api.Trace) api.Summary {
	traceByID := map[string][]api.Trace{}

	for _, trace := range traces {
		traces := traceByID[trace.ID()]
		if traces == nil {
			traces = []api.Trace{}
		}
		traceByID[trace.ID()] = append(traces, trace)
	}

	return exec.NewSummary(traceByID)
}
