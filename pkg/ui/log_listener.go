package ui

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sha1n/bert/api"
)

// LoggingProgressListener logs progress events with the contextual information a logger.
type LoggingProgressListener struct {
}

// NewLoggingProgressListener creates a new listener which logs to the standard logger.
func NewLoggingProgressListener() api.Listener {
	return LoggingProgressListener{}
}

// OnBenchmarkStart logs an info message
func (l LoggingProgressListener) OnBenchmarkStart() {
	slog.Info("Starting benchmark...")
}

// OnBenchmarkEnd logs an info message
func (l LoggingProgressListener) OnBenchmarkEnd() {
	slog.Info("Benchmark finished")
}

// OnScenarioStart logs an info message with the specified ID
func (l LoggingProgressListener) OnScenarioStart(id api.ID) {
	slog.Info(fmt.Sprintf("[%s] starting...", yellow.Sprint(id)))
}

// OnScenarioEnd logs an info message with the specified ID
func (l LoggingProgressListener) OnScenarioEnd(id api.ID) {
	slog.Info(fmt.Sprintf("[%s] finished", yellow.Sprint(id)))
}

// OnError logs an error message with the specified ID and error details
func (l LoggingProgressListener) OnError(id api.ID, err error) {
	slog.Error(fmt.Sprintf("[%s] error: %v", yellow.Sprint(id), err))
}

// OnMessage logs an info message with the specified ID and text
func (l LoggingProgressListener) OnMessage(id api.ID, message string) {
	slog.Info(fmt.Sprintf("[%s] %s", yellow.Sprint(id), message))
}

// OnMessagef logs an info message with the specified ID and formatted message text
func (l LoggingProgressListener) OnMessagef(id api.ID, format string, args ...interface{}) {
	if slog.Default().Enabled(context.Background(), slog.LevelInfo) {
		slog.Info(fmt.Sprintf("[%s] %s", yellow.Sprint(id), fmt.Sprintf(format, args...)))
	}
}
