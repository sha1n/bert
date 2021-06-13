package pkg

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/sha1n/benchy/api"
	clibtest "github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestMinimalProgressViewOutput(t *testing.T) {
	ctx := api.NewIOContext()
	ctx.Tty = true
	ctx.StdoutWriter = new(bytes.Buffer)
	ctx.StderrWriter = new(bytes.Buffer)
	spec := aBasicSpecWith(true, 2)
	scenarioID := spec.Scenarios[0].ID()

	expectedError := errors.New(clibtest.RandomString())

	progView := NewMinimalProgressView(spec, fakeTermDimensions, ctx).(*MinimalProgressView)

	progView.OnBenchmarkStart()

	// round one
	progView.OnScenarioStart(scenarioID)
	progView.OnError(scenarioID, expectedError)
	assert.Error(t, progView.progressInfoByID[scenarioID].lastError)
	assert.Equal(t, expectedError, progView.progressInfoByID[scenarioID].lastError)

	time.Sleep(time.Nanosecond) // make sure mean is not zero

	progView.OnScenarioEnd(scenarioID)
	assert.True(t, progView.progressInfoByID[scenarioID].mean > 0)

	// round two
	progView.OnScenarioStart(scenarioID)
	assert.Equal(t, 1, progView.progressInfoByID[scenarioID].executions)

	progView.OnScenarioEnd(scenarioID)
	assert.True(t, progView.progressInfoByID[scenarioID].mean > 0)

	progView.OnBenchmarkEnd()
}

func TestMinimalProgressViewEndNotStartedStateContract(t *testing.T) {
	testProgressViewEndNotStartedStateContract(
		t,
		NewMinimalProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestMinimalProgressViewEndStateContract(t *testing.T) {
	testProgressViewEndStateContract(
		t,
		NewMinimalProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestMinimalProgressViewStartAlreadyEndedStateContract(t *testing.T) {
	testProgressViewStartAlreadyEndedStateContract(
		t,
		NewMinimalProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func TestMinimalProgressViewStartStateContract(t *testing.T) {
	testProgressViewStartStateContract(
		t,
		NewMinimalProgressView(aBasicSpecWith(true, 1), fakeTermDimensions, api.NewIOContext()),
	)
}

func testProgressViewEndNotStartedStateContract(t *testing.T, progView api.Listener) {
	assert.Panics(t, func() {
		progView.OnBenchmarkEnd()
	})
}

func testProgressViewEndStateContract(t *testing.T, progView api.Listener) {
	progView.OnBenchmarkStart()
	progView.OnBenchmarkEnd()

	assert.Panics(t, func() {
		progView.OnBenchmarkEnd()
	})
}

func testProgressViewStartAlreadyEndedStateContract(t *testing.T, progView api.Listener) {
	progView.OnBenchmarkStart()
	progView.OnBenchmarkEnd()

	assert.Panics(t, func() {
		progView.OnBenchmarkStart()
	})
}

func testProgressViewStartStateContract(t *testing.T, progView api.Listener) {
	progView.OnBenchmarkStart()
	assert.Panics(t, func() {
		progView.OnBenchmarkStart()
	})
}
