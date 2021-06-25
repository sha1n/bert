package pkg

import (
	"time"

	"github.com/sha1n/bert/api"
)

// NewFakeTrace creates a fake trace with the specified data
func NewFakeTrace(id string, elapsed, userTime, sysTime time.Duration, err error) api.Trace {
	return &trace{
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

	return NewSummary(traceByID)
}
