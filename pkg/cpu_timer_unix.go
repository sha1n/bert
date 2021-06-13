// +build linux darwin
// +build amd64 arm64

package pkg

import (
	"syscall"
	"time"

	"github.com/sha1n/benchy/api"
)

type unixChildrenCPUTimer struct {
	r                  syscall.Rusage
	who                int
	perceivedStartTime time.Time
	sysTimeStart       time.Time
	usrTimeStart       time.Time
}

// NewChildrenCPUTimer returns a new CPUTimer that measures sub-processes CPU time using system calls.
func NewChildrenCPUTimer() api.CPUTimer {
	return newCPUTimer(syscall.RUSAGE_CHILDREN)
}

// NewSelfCPUTimer returns a new CPUTimer that measures this process' CPU time using system calls.
func NewSelfCPUTimer() api.CPUTimer {
	return newCPUTimer(syscall.RUSAGE_SELF)
}

func newCPUTimer(who int) api.CPUTimer {
	return &unixChildrenCPUTimer{
		who: who,
	}
}

func (t *unixChildrenCPUTimer) Start() api.ElapsedCPUTimeFn {
	err := syscall.Getrusage(t.who, &t.r)
	if err != nil {
		panic(err)
	}

	t.perceivedStartTime = time.Now()
	t.usrTimeStart = time.Unix(t.r.Utime.Unix())
	t.sysTimeStart = time.Unix(t.r.Stime.Unix())

	return t.Elapsed
}

func (t *unixChildrenCPUTimer) Elapsed() (perceived time.Duration, usr time.Duration, sys time.Duration) {
	err := syscall.Getrusage(t.who, &t.r)
	if err != nil {
		panic(err)
	}

	perceivedTime := time.Now().Sub(t.perceivedStartTime)
	usrTime := time.Unix(t.r.Utime.Unix()).Sub(t.usrTimeStart)
	sysTime := time.Unix(t.r.Stime.Unix()).Sub(t.sysTimeStart)

	return perceivedTime, usrTime, sysTime
}
