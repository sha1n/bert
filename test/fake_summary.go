package test

import (
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
)

type fakeTrace struct {
	id      string
	elapsed time.Duration
	error   error
}

// NewFakeTrace creates a fake trace with the specified data
func NewFakeTrace(id string, elapsed time.Duration, err error) api.Trace {
	return &fakeTrace{
		id:      id,
		elapsed: elapsed,
		error:   err,
	}
}

func (t *fakeTrace) ID() string {
	return t.id
}

func (t *fakeTrace) Elapsed() time.Duration {
	return t.elapsed
}

func (t *fakeTrace) Error() error {
	return t.error
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

	return pkg.NewSummary(traceByID)
}
