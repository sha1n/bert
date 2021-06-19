package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/sha1n/gommons/pkg/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogProgressListener(t *testing.T) {
	expected := LoggingProgressListener{logger: log.StandardLogger()}

	assert.Equal(t, expected, NewLoggingProgressListener())
}

func TestLogProgressListener_OnBenchmarkStart(t *testing.T) {
	l := newInterceptableLogProgressListener()

	l.OnBenchmarkStart()

	assert.NotEmpty(t, getLoggedString(l.logger))
}

func TestLogProgressListener_OnBenchmarkEnd(t *testing.T) {
	l := newInterceptableLogProgressListener()

	l.OnBenchmarkEnd()

	assert.NotEmpty(t, getLoggedString(l.logger))
}

func TestLogProgressListener_OnScenarioStart(t *testing.T) {
	l := newInterceptableLogProgressListener()
	expectedScenarioID := test.RandomString()

	l.OnScenarioStart(expectedScenarioID)

	assert.Contains(t, getLoggedString(l.logger), expectedScenarioID)
}

func TestLogProgressListener_OnScenarioEnd(t *testing.T) {
	l := newInterceptableLogProgressListener()
	expectedScenarioID := test.RandomString()

	l.OnScenarioEnd(expectedScenarioID)

	assert.Contains(t, getLoggedString(l.logger), expectedScenarioID)
}

func TestLogProgressListener_OnError(t *testing.T) {
	l := newInterceptableLogProgressListener()
	expectedScenarioID := test.RandomString()
	expectedError := errors.New(test.RandomString())

	l.OnError(expectedScenarioID, expectedError)

	actual := getLoggedString(l.logger)
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, expectedError.Error())
}

func TestLogProgressListener_OnMessage(t *testing.T) {
	l := newInterceptableLogProgressListener()
	expectedScenarioID := test.RandomString()
	expectedMessage := test.RandomString()

	l.OnMessage(expectedScenarioID, expectedMessage)

	actual := getLoggedString(l.logger)
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, expectedMessage)
}

func TestLogProgressListener_OnMessagef(t *testing.T) {
	l := newInterceptableLogProgressListener()
	expectedScenarioID := test.RandomString()
	expectedMessageFormat := "format: %s"
	expectedMessageParam := test.RandomString()

	l.OnMessagef(expectedScenarioID, expectedMessageFormat, expectedMessageParam)

	actual := getLoggedString(l.logger)
	assert.Contains(t, actual, expectedScenarioID)
	assert.Contains(t, actual, fmt.Sprintf(expectedMessageFormat, expectedMessageParam))
}

func newInterceptableLogProgressListener() LoggingProgressListener {
	logger := newInterceptingLogger()
	l := LoggingProgressListener{
		logger: logger,
	}

	return l
}

func newInterceptingLogger() *log.Logger {
	logger := log.New()
	buffer := new(bytes.Buffer)
	logger.Out = buffer

	return logger
}

func getLoggedString(logger *log.Logger) string {
	return logger.Out.(*bytes.Buffer).String()
}
