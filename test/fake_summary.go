package test

import (
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
)

type fake_trace struct {
	id      string
	elapsed time.Duration
	error   error
}

func NewFakeTrace(id string, elapsed time.Duration, err error) api.Trace {
	return &fake_trace{
		id:      id,
		elapsed: elapsed,
		error:   err,
	}
}

func (t *fake_trace) ID() string {
	return t.id
}

func (t *fake_trace) Elapsed() time.Duration {
	return t.elapsed
}

func (t *fake_trace) Error() error {
	return t.error
}

func NewFakeSummary(traces ...api.Trace) api.Summary {
	traceByID := map[string][]api.Trace{}

	for i := range traces {
		trace := traces[i]
		traces := traceByID[trace.ID()]
		if traces == nil {
			traces = []api.Trace{}
		}
		traceByID[trace.ID()] = append(traces, trace)
	}

	return pkg.NewSummary(traceByID)
}
