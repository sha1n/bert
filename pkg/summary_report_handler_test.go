package pkg

import (
	"errors"
	"testing"

	"github.com/sha1n/benchy/api"

	"github.com/stretchr/testify/assert"
)

func TestSummaryReportHandler(t *testing.T) {
	testSummaryReportHandlerFinalizeWith(t, nil)
}

func TestSummaryReportHandlerFinalizeError(t *testing.T) {
	testSummaryReportHandlerFinalizeWith(t, errors.New("test error"))
}

func testSummaryReportHandlerFinalizeWith(t *testing.T, expectedError error) {
	tracer := NewTracer(1)
	expectedSpec := exampleSpec()
	expectedCtx := api.ReportContext{}

	interceptor := newWriteReportInterceptor(expectedError)

	handler := NewSummaryReportHandler(expectedSpec, expectedCtx, interceptor.intercept)

	handler.Subscribe(tracer.Stream())

	// Fire one trace event
	tracer.Start(expectedSpec.Scenarios[0])(nil)

	actualError := handler.Finalize()

	assert.Equal(t, expectedError, actualError)

	assert.Equal(t, expectedSpec, interceptor.capturedSpec)
	assert.Equal(t, expectedCtx, interceptor.capturedCtx)
	assert.Equal(t, 1, len(interceptor.capturedSummary.IDs()))
}

type writeReportInterceptor struct {
	capturedSpec    api.BenchmarkSpec
	capturedCtx     api.ReportContext
	capturedSummary api.Summary
	expectedError   error
}

func newWriteReportInterceptor(expectedError error) *writeReportInterceptor {
	return &writeReportInterceptor{
		expectedError: expectedError,
	}
}

func (i *writeReportInterceptor) intercept(summary api.Summary, spec api.BenchmarkSpec, ctx api.ReportContext) error {
	i.capturedSummary = summary
	i.capturedSpec = spec
	i.capturedCtx = ctx

	return i.expectedError
}

func exampleSpec() api.BenchmarkSpec {
	return api.BenchmarkSpec{
		Executions: 1,
		Scenarios: []api.ScenarioSpec{
			{
				Name: "scenario",
				Command: &api.CommandSpec{
					Cmd: []string{"cmd"},
				},
			},
		},
	}
}
