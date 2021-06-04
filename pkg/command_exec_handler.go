package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
)

type RunCommandFn = func(cmd *exec.Cmd, ctx api.IOContext) (err error)

func RunCommandFnFor(ce *commandExecutor) RunCommandFn {
	// if ce.pipeStdout || ce.pipeStderr || !ce.ctx.Tty || ce.ctx.DisbaleRichTerminalEffects {
	// 	return RunCommand
	// }
	return RunCommandWithProgressIndicator
}

func RunCommand(cmd *exec.Cmd, ctx api.IOContext) (err error) {
	registerInterruptGuard(cmd, func(c *exec.Cmd, s os.Signal) {
		onShutdownSignal(c, s)
	})

	return cmd.Run()
}

func RunCommandWithProgressIndicator(cmd *exec.Cmd, ctx api.IOContext) (err error) {
	cursor := termite.NewCursor(termite.StdoutWriter)
	cursor.Hide()
	defer cursor.Show()

	registerInterruptGuard(cmd, func(c *exec.Cmd, s os.Signal) {
		cursor.Show()
		onShutdownSignal(c, s)
	})

	spinner := termite.NewSpinner(ctx.StdoutWriter, "Preparing...", time.Millisecond*100, &spinnerFormatter{})

	if _, err = spinner.Start(); err == nil {
		_ = spinner.SetTitle(fmt.Sprintf("Executing command %v", cmd.Args))

		err = cmd.Run()
		if err != nil {
			spinner.SetTitle(fmt.Sprintf("Command failed! Error: %v", err))
		}
	}

	spinner.Stop("")

	return err
}

func onShutdownSignal(execCmd *exec.Cmd, sig os.Signal) {
	if sig == os.Interrupt {
		log.Debugf("Got %s signal. Forwarding to %s...", sig, execCmd.Args[0])
		execCmd.Process.Signal(sig)

		os.Exit(1)
	}
}

// channel is returned for testing...
func registerInterruptGuard(execCmd *exec.Cmd, handleFn func(*exec.Cmd, os.Signal)) (context.CancelFunc, chan os.Signal) {
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
				handleFn(execCmd, sig)
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

type spinnerFormatter struct{}

var cyan = color.New(color.FgCyan)

// FormatTitle returns the input title as is
func (f *spinnerFormatter) FormatTitle(s string) string {
	return s
}

// FormatIndicator returns the input char as is
func (f *spinnerFormatter) FormatIndicator(char string) string {
	if log.StandardLogger().Level != log.PanicLevel {
		return cyan.Sprintf("%s%s", strings.Repeat(" ", 3), char)
	}
	return cyan.Sprint(char)
}

// CharSeq returns the default character sequence.
func (f *spinnerFormatter) CharSeq() []string {
	return termite.DefaultSpinnerCharSeq()
}
