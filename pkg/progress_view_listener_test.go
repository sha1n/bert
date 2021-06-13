package pkg

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/benchy/api"
	clibtest "github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

// ProgressView is very tricky low-roi API to test.
// This test makes sure that nothing gets stuck and verifies that data is written to the right stream.
func TestProgressViewOutput(t *testing.T) {
	ctx := api.NewIOContext()
	ctx.Tty = true
	ctx.StdoutWriter = new(bytes.Buffer)
	ctx.StderrWriter = new(bytes.Buffer)
	spec := aBasicSpecWith(true, 2)
	scenarioID := spec.Scenarios[0].ID()

	errorMessage := clibtest.RandomString()

	progView := NewProgressView(spec, fakeTermDimensions, ctx).(*ProgressView)

	progView.OnBenchmarkStart()

	// round one
	progView.OnScenarioStart(scenarioID)
	progView.OnError(scenarioID, errors.New(errorMessage))
	stdoutEventuallyContains(t, errorMessage, ctx)

	time.Sleep(time.Nanosecond) // make sure mean is not zero

	progView.OnScenarioEnd(scenarioID)
	assert.True(t, progView.progressInfoByID[scenarioID].mean > 0)

	// round two
	progView.OnScenarioStart(scenarioID)
	assert.Equal(t, 1, progView.progressInfoByID[scenarioID].executions)

	progView.OnScenarioEnd(scenarioID)
	assert.False(t, progView.progressInfoByID[scenarioID].tick(""), "progress bar is expected to finish")
	assert.True(t, progView.progressInfoByID[scenarioID].mean > 0)

	progView.OnBenchmarkEnd()
}

func TestProgressViewStartStateContract(t *testing.T) {
	testProgressViewStartStateContract(
		t,
		NewProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestProgressViewStartAlreadyEndedStateContract(t *testing.T) {
	testProgressViewStartAlreadyEndedStateContract(
		t,
		NewProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestProgressViewEndStateContract(t *testing.T) {
	testProgressViewEndStateContract(
		t,
		NewProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestProgressViewEndNotStartedStateContract(t *testing.T) {
	testProgressViewEndNotStartedStateContract(
		t,
		NewProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func Test_progressInfo_calculateNewApproxMean(t *testing.T) {
	type fields struct {
		executions         int
		expectedExecutions int
		mean               time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		elapsed time.Duration
		want    time.Duration
	}{
		{name: "no executions", fields: fields{executions: 0, expectedExecutions: 10, mean: 0}, elapsed: time.Millisecond * 1, want: time.Millisecond * 1},
		{name: "2 out of 10", fields: fields{executions: 1, expectedExecutions: 10, mean: time.Millisecond * 1}, elapsed: time.Millisecond * 3, want: time.Millisecond * 2},
		{name: "10 out of 10 - DONE!", fields: fields{executions: 10, expectedExecutions: 10, mean: time.Millisecond * 1}, elapsed: time.Millisecond * 100, want: time.Millisecond * 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := progressInfo{
				minimalProgressInfo: minimalProgressInfo{
					executions:         tt.fields.executions,
					expectedExecutions: tt.fields.expectedExecutions,
					mean:               tt.fields.mean,
				},
			}
			if got := pi.calculateNewApproxMean(tt.elapsed); got != tt.want {
				t.Errorf("progressInfo.calculateNewApproxMean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_progressInfo_calculateETA(t *testing.T) {
	type fields struct {
		executions         int
		expectedExecutions int
		mean               time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{name: "no executions", fields: fields{executions: 0, expectedExecutions: 10, mean: 0}, want: 0},
		{name: "10% done", fields: fields{executions: 1, expectedExecutions: 10, mean: time.Millisecond * 1}, want: time.Millisecond * 9},
		{name: "done", fields: fields{executions: 10, expectedExecutions: 10, mean: time.Millisecond * 1}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := progressInfo{
				minimalProgressInfo: minimalProgressInfo{
					executions:         tt.fields.executions,
					expectedExecutions: tt.fields.expectedExecutions,
					mean:               tt.fields.mean,
				},
			}
			if got := pi.calculateETA(); got != tt.want {
				t.Errorf("progressInfo.calculateETA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stdoutEventuallyContains(t *testing.T, want string, ctx api.IOContext) {
	assert.Eventually(t, stdoutContainsFn(ctx, want), time.Second*30, time.Millisecond*10)
}

func stdoutContainsFn(ctx api.IOContext, want string) func() bool {
	return func() bool {
		return strings.Contains(ctx.StdoutWriter.(*bytes.Buffer).String(), want)
	}
}

func fakeTermDimensions() (int, int) {
	return 100, 100
}
