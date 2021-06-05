package pkg

import (
	"fmt"

	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

type LogProgressListener struct {
	logger *log.Logger
}

func NewLogProgressListener() api.Listener {
	return LogProgressListener{
		logger: log.StandardLogger(),
	}
}

func (l LogProgressListener) OnBenchmarkStart() {
	l.logger.Info("Starting benchmark...")
}

func (l LogProgressListener) OnBenchmarkEnd() {
	l.logger.Info("Benchmark finished")
}

func (l LogProgressListener) OnScenarioStart(id api.ID) {
	l.logger.Infof("[%s] Starting...", id)
}

func (l LogProgressListener) OnScenarioEnd(id api.ID) {
	l.logger.Infof("[%s] finished", id)
}

func (l LogProgressListener) OnError(id api.ID, err error) {
	l.logger.Errorf("[%s] Error: %v", id, err)
}

func (l LogProgressListener) OnMessage(id api.ID, message string) {
	l.logger.Infof("[%s] %s", id, message)
}

func (l LogProgressListener) OnMessagef(id api.ID, format string, args ...interface{}) {
	if l.logger.IsLevelEnabled(log.InfoLevel) {
		l.logger.Infof("[%s] %s", id, fmt.Sprintf(format, args...))
	}
	
}