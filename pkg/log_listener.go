package pkg

import (
	"fmt"

	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// LoggingProgressListener logs progress events with the contextual information a logger.
type LoggingProgressListener struct {
	logger *log.Logger
}

// NewLoggingProgressListener creates a new listener which logs to the standard logger.
func NewLoggingProgressListener() api.Listener {
	return LoggingProgressListener{
		logger: log.StandardLogger(),
	}
}

// OnBenchmarkStart logs an info message
func (l LoggingProgressListener) OnBenchmarkStart() {
	l.logger.Info("Starting benchmark...")
}

// OnBenchmarkEnd logs an info message
func (l LoggingProgressListener) OnBenchmarkEnd() {
	l.logger.Info("Benchmark finished")
}

// OnScenarioStart logs an info message with the specified ID
func (l LoggingProgressListener) OnScenarioStart(id api.ID) {
	l.logger.Infof("[%s] starting...", yellow.Sprint(id))
}

// OnScenarioEnd logs an info message with the specified ID
func (l LoggingProgressListener) OnScenarioEnd(id api.ID) {
	l.logger.Infof("[%s] finished", yellow.Sprint(id))
}

// OnError logs an error message with the specified ID and error details
func (l LoggingProgressListener) OnError(id api.ID, err error) {
	l.logger.Errorf("[%s] error: %v", yellow.Sprint(id), err)
}

// OnMessage logs an info message with the specified ID and text
func (l LoggingProgressListener) OnMessage(id api.ID, message string) {
	l.logger.Infof("[%s] %s", yellow.Sprint(id), message)
}

// OnMessagef logs an info message with the specified ID and formatted message text
func (l LoggingProgressListener) OnMessagef(id api.ID, format string, args ...interface{}) {
	if l.logger.IsLevelEnabled(log.InfoLevel) {
		l.logger.Infof("[%s] %s", yellow.Sprint(id), fmt.Sprintf(format, args...))
	}

}
