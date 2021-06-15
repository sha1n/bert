package pkg

import (
	"time"

	"github.com/sha1n/bert/api"
)

// PerceivedTimeCPUTimer NOOP implementation of CPUTimer
type PerceivedTimeCPUTimer struct {
	perceivedStartTime time.Time
}

// Start ...
func (t *PerceivedTimeCPUTimer) Start() api.ElapsedCPUTimeFn {
	t.perceivedStartTime = time.Now()

	return t.Elapsed
}

// Elapsed return (measured perceived time, 0, 0)
func (t *PerceivedTimeCPUTimer) Elapsed() (perceived time.Duration, usr time.Duration, sys time.Duration) {
	perceivedTime := time.Now().Sub(t.perceivedStartTime)

	return perceivedTime, time.Nanosecond * 0, time.Nanosecond * 0
}
