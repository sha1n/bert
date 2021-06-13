package api

import "time"

// ElapsedCPUTimeFn returns a pair of CPU time measurements, one for user time and one for system time.
type ElapsedCPUTimeFn = func() (perceived time.Duration, user time.Duration, system time.Duration)

// CPUTimer an abstraction for a platform system CPU timer
type CPUTimer interface {
	Start() ElapsedCPUTimeFn
	Elapsed() (perceived time.Duration, usr time.Duration, sys time.Duration)
}
