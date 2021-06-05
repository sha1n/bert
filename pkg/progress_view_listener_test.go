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
	spec := aBasicSpecWith(true, 1)
	scenarioID := spec.Scenarios[0].ID()

	message1, message2, errorMessage := clibtest.RandomString(), clibtest.RandomString(), clibtest.RandomString()

	progView := NewProgressView(spec, fakeTermWidth, ctx)
	progView.OnBenchmarkStart()

	progView.OnScenarioStart(scenarioID)
	stdoutEventuallyContains(t, scenarioID, ctx)

	progView.OnMessage(scenarioID, message1)
	stdoutEventuallyContains(t, message1, ctx)

	progView.OnMessagef(scenarioID, "%s", message2)
	stdoutEventuallyContains(t, message2, ctx)

	progView.OnError(scenarioID, errors.New(errorMessage))
	stdoutEventuallyContains(t, errorMessage, ctx)

	progView.OnScenarioEnd(scenarioID)
	progView.OnBenchmarkEnd()
}

func TestProgressViewStartStateContract(t *testing.T) {
	progView := NewProgressView(aBasicSpecWith(true, 1), fakeTermWidth, api.NewIOContext())

	progView.OnBenchmarkStart()
	assert.Panics(t, func() {
		progView.OnBenchmarkStart()
	})
}

func TestProgressViewStartAlreadyEndedStateContract(t *testing.T) {
	progView := NewProgressView(aBasicSpecWith(true, 1), fakeTermWidth, api.NewIOContext())
	progView.OnBenchmarkStart()
	progView.OnBenchmarkEnd()

	assert.Panics(t, func() {
		progView.OnBenchmarkStart()
	})
}

func TestProgressViewEndStateContract(t *testing.T) {
	progView := NewProgressView(aBasicSpecWith(true, 1), fakeTermWidth, api.NewIOContext())
	progView.OnBenchmarkStart()
	progView.OnBenchmarkEnd()

	assert.Panics(t, func() {
		progView.OnBenchmarkEnd()
	})
}

func TestProgressViewEndNotStartedStateContract(t *testing.T) {
	progView := NewProgressView(aBasicSpecWith(true, 1), fakeTermWidth, api.NewIOContext())

	assert.Panics(t, func() {
		progView.OnBenchmarkEnd()
	})
}

func stdoutEventuallyContains(t *testing.T, want string, ctx api.IOContext) {
	assert.Eventually(t, stdoutContainsFn(ctx, want), time.Second*30, time.Millisecond*10)
}

func stdoutContainsFn(ctx api.IOContext, want string) func() bool {
	return func() bool {
		return strings.Contains(ctx.StdoutWriter.(*bytes.Buffer).String(), want)
	}
}

func fakeTermWidth() int {
	return 100
}
