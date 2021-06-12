// +build windows
// +build amd64 arm

package pkg

import (
	"time"
)

// NoopCPUTimer NOOP implementation of CPUTimer
type NoopCPUTimer struct{}

func newChildrenCPUTimer() CPUTimer {
	return NoopCPUTimer{}
}

// Start ...
func (t NoopCPUTimer) Start() func() (time.Duration, time.Duration) {
	return t.Elapsed
}

// Elapsed return 0, 0
func (t NoopCPUTimer) Elapsed() (usr time.Duration, sys time.Duration) {
	return time.Nanosecond * 0, time.Nanosecond * 0
}
