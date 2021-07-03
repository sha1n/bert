package report_handlers

import (
	"errors"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg/exec"

	"github.com/stretchr/testify/assert"
)

func TestStreamReportHandler(t *testing.T) {
	testStreamReportHandlerFinalizeWith(t, nil, nil)
}

func TestStreamReportHandlerFinalizeHandleError(t *testing.T) {
	testStreamReportHandlerFinalizeWith(t, errors.New("test error"), nil)
}

func TestStreamReportHandlerFinalizeTraceError(t *testing.T) {
	testStreamReportHandlerFinalizeWith(t, nil, errors.New("test error"))
}

func testStreamReportHandlerFinalizeWith(t *testing.T, expectedHandleError error, expectedTraceError error) {
	tracer := exec.NewTracer(1)
	expectedSpec := exampleSpec()
	expectedCtx := api.ReportContext{}

	interceptor := newHandleInterceptor(expectedHandleError)

	handler := NewStreamReportHandler(expectedSpec, expectedCtx, interceptor.intercept)

	handler.Subscribe(tracer.Stream())

	// Fire one trace event
	tracer.Start(expectedSpec.Scenarios[0])(&api.ExecutionInfo{}, expectedTraceError)

	assert.NoError(t, handler.Finalize())

	assert.Equal(t, expectedHandleError, interceptor.expectedError)
	assert.Equal(t, expectedSpec.Scenarios[0].ID(), interceptor.capturedTrace.ID())
	assert.Equal(t, expectedTraceError, interceptor.capturedTrace.Error())
}

type handleInterceptor struct {
	capturedTrace api.Trace
	expectedError error
}

func newHandleInterceptor(expectedError error) *handleInterceptor {
	return &handleInterceptor{
		expectedError: expectedError,
	}
}

func (i *handleInterceptor) intercept(trace api.Trace) error {
	i.capturedTrace = trace

	return i.expectedError
}
