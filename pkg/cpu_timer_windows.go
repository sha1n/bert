// +build windows
// +build amd64 arm

package pkg

import (
	"time"

	"github.com/sha1n/bert/api"
)

// NoopCPUTimer NOOP implementation of CPUTimer
type NoopCPUTimer struct {
	perceivedStartTime time.Time
}

// NewChildrenCPUTimer returns a NOOP CPUTimer implementation.
func NewChildrenCPUTimer() api.CPUTimer {
	return NoopCPUTimer{}
}

// NewSelfCPUTimer returns a NOOP CPUTimer implementation.
func NewSelfCPUTimer() api.CPUTimer {
	return NoopCPUTimer{}
}

// Start ...
func (t NoopCPUTimer) Start() api.ElapsedCPUTimeFn {
	t.perceivedStartTime = time.Now()

	return t.Elapsed
}

// Elapsed return 0, 0
func (t NoopCPUTimer) Elapsed() (perceived time.Duration, usr time.Duration, sys time.Duration) {
	perceivedTime := time.Now().Sub(t.perceivedStartTime)

	return perceivedTime, time.Nanosecond * 0, time.Nanosecond * 0
}
