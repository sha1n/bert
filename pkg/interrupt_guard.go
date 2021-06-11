package pkg

import (
	"context"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
)

// RegisterInterruptGuard registers an interrupt signal handler and returns a cancellation function.
func RegisterInterruptGuard(handleFn func(os.Signal)) context.CancelFunc {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	startWG := &sync.WaitGroup{}
	startWG.Add(1)

	go func() {
		startWG.Done()

		select {
		case sig, ok := <-c:
			if ok {
				handleFn(sig)
			}

		case <-ctx.Done():
			signal.Stop(c)

			close(c)
			log.Debug("Context cancelled - OK!")
		}
	}()

	startWG.Wait()

	return cancel
}

// channel is returned for testing...
func registerInterruptGuard(handleFn func(os.Signal)) (context.CancelFunc, chan os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	startWG := &sync.WaitGroup{}
	startWG.Add(1)

	go func() {
		startWG.Done()

		select {
		case sig, ok := <-c:
			if ok {
				handleFn(sig)
			}

		case <-ctx.Done():
			signal.Stop(c)

			close(c)
			log.Debug("Context cancelled - OK!")
		}
	}()

	startWG.Wait()

	return cancel, c
}
