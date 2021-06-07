package pkg

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterInterruptGuard(t *testing.T) {
	call := make(chan bool)
	_, c := registerInterruptGuard(func(s os.Signal) {
		call <- true
	})

	c <- os.Interrupt
	assert.Eventually(t, func() bool { return <-call }, time.Second*10, time.Millisecond)
}

func TestRegisterInterruptGuardCancellation(t *testing.T) {
	cancel, c := registerInterruptGuard(func(s os.Signal) {})
	cancel()

	assert.Panics(t, func() {
		c <- os.Interrupt // this should fail panic because the channel is closed
	})

}
