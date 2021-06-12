// +build linux darwin
// +build amd64 arm64

package pkg

import (
	"syscall"
	"time"
)

type unixChildrenCPUTimer struct {
	r            syscall.Rusage
	sysTimeStart time.Time
	usrTimeStart time.Time
}

func newChildrenCPUTimer() CPUTimer {
	return &unixChildrenCPUTimer{}
}

func (t *unixChildrenCPUTimer) Start() func() (time.Duration, time.Duration) {
	err := syscall.Getrusage(syscall.RUSAGE_CHILDREN, &t.r)
	if err != nil {
		panic(err)
	}

	t.sysTimeStart = time.Unix(t.r.Stime.Unix())
	t.usrTimeStart = time.Unix(t.r.Utime.Unix())

	return t.Elapsed
}

func (t *unixChildrenCPUTimer) Elapsed() (usr time.Duration, sys time.Duration) {
	err := syscall.Getrusage(syscall.RUSAGE_CHILDREN, &t.r)
	if err != nil {
		panic(err)
	}

	sysTimeEnd := time.Unix(t.r.Stime.Unix())
	usrTimeEnd := time.Unix(t.r.Utime.Unix())

	return usrTimeEnd.Sub(t.usrTimeStart), sysTimeEnd.Sub(t.sysTimeStart)
}
