// +build linux darwin
// +build amd64 arm64

package pkg

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewChildrenCPUTimer(t *testing.T) {
	timer := NewChildrenCPUTimer()

	assert.Equal(t, syscall.RUSAGE_CHILDREN, timer.(*unixChildrenCPUTimer).who)
}

func TestNewSelfCPUTimer(t *testing.T) {
	timer := NewSelfCPUTimer()

	assert.Equal(t, syscall.RUSAGE_SELF, timer.(*unixChildrenCPUTimer).who)
}

func Test_SelfCPUTimer_Elapsed_Unloaded(t *testing.T) {
	timer := NewSelfCPUTimer()

	elapsed := timer.Start()
	perceived, usr, sys := elapsed()

	assert.GreaterOrEqual(t, time.Millisecond*1, perceived)
	assert.GreaterOrEqual(t, time.Millisecond*1, usr)
	assert.GreaterOrEqual(t, time.Millisecond*1, sys)
}

func Test_SelfCPUTimer_Elapsed_Loaded(t *testing.T) {
	timer := NewSelfCPUTimer()

	elapsed := timer.Start()

	assert.Eventually(t,
		func() bool {
			perceived, usr, sys := elapsed()

			return perceived > time.Millisecond*1 && usr > time.Millisecond*1 && sys > time.Millisecond*1
		},
		time.Second*1,
		time.Nanosecond*1,
	)
}

func Test_ChildrenCPUTimer_Elapsed_Loaded(t *testing.T) {
	timer := NewChildrenCPUTimer()

	elapsed := timer.Start()
	exec.Command("go", "list", "./...").Run()
	perceived, usr, sys := elapsed()

	assert.GreaterOrEqual(t, perceived, time.Nanosecond*1)
	assert.GreaterOrEqual(t, usr, time.Nanosecond*1)
	assert.GreaterOrEqual(t, sys, time.Nanosecond*1)
}

func Test_ChildrenCPUTimer_Elapsed_Unloaded(t *testing.T) {
	timer := NewChildrenCPUTimer()

	elapsed := timer.Start()
	perceived, usr, sys := elapsed()

	assert.GreaterOrEqual(t, perceived, time.Nanosecond*0)
	assert.Equal(t, time.Nanosecond*0, usr)
	assert.Equal(t, time.Nanosecond*0, sys)
}
