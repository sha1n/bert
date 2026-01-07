package ui

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestNewLogProgressListener(t *testing.T) {
	assert.Equal(t, LoggingProgressListener{}, NewLoggingProgressListener())
}

func TestLogProgressListener_OnBenchmarkStart(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()

	l.OnBenchmarkStart()

	assert.NotEmpty(t, buf.String())
}

func TestLogProgressListener_OnBenchmarkEnd(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()

	l.OnBenchmarkEnd()

	assert.NotEmpty(t, buf.String())
}

func TestLogProgressListener_OnScenarioStart(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()
	expectedScenarioID := test.RandomString()

	l.OnScenarioStart(expectedScenarioID)

	assert.Contains(t, buf.String(), expectedScenarioID)
}

func TestLogProgressListener_OnScenarioEnd(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()
	expectedScenarioID := test.RandomString()

	l.OnScenarioEnd(expectedScenarioID)

	assert.Contains(t, buf.String(), expectedScenarioID)
}

func TestLogProgressListener_OnError(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()
	expectedScenarioID := test.RandomString()
	expectedError := errors.New(test.RandomString())

	l.OnError(expectedScenarioID, expectedError)

	actual := buf.String()
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, expectedError.Error())
}

func TestLogProgressListener_OnMessage(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()
	expectedScenarioID := test.RandomString()
	expectedMessage := test.RandomString()

	l.OnMessage(expectedScenarioID, expectedMessage)

	actual := buf.String()
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, expectedMessage)
}

func TestLogProgressListener_OnMessagef(t *testing.T) {
	l := NewLoggingProgressListener()
	buf, restore := interceptSlog()
	defer restore()
	expectedScenarioID := test.RandomString()
	expectedMessageFormat := "format: %s"
	expectedMessageParam := test.RandomString()

	l.OnMessagef(expectedScenarioID, expectedMessageFormat, expectedMessageParam)

	actual := buf.String()
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, fmt.Sprintf(expectedMessageFormat, expectedMessageParam))
}

func interceptSlog() (*bytes.Buffer, func()) {
	original := slog.Default()
	buf := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))

	return buf, func() { slog.SetDefault(original) }
}
