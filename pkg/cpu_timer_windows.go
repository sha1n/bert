// +build windows
// +build amd64 arm

package pkg

import (
	"time"

	"github.com/sha1n/bert/api"
)

// NewChildrenCPUTimer returns a NOOP CPUTimer implementation.
func NewChildrenCPUTimer() api.CPUTimer {
	return PerceivedTimeCPUTimer{}
}

// NewSelfCPUTimer returns a NOOP CPUTimer implementation.
func NewSelfCPUTimer() api.CPUTimer {
	return PerceivedTimeCPUTimer{}
}
