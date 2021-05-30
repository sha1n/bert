package cli

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	printfRed   = color.New(color.FgRed).Printf
	printRed    = color.New(color.FgRed).Print
	sprintRed   = color.New(color.FgRed).Sprint
	sprintGreen = color.New(color.FgGreen).Sprint
	sprintBold  = color.New(color.Bold).Sprint
)

// this writer should be used when background threads write messages without a new line suffix to stdout.
// this is common with progress indicartors from the termite package.
type alwaysRewritingWriter struct {
	writer io.Writer
}

func (sw *alwaysRewritingWriter) Write(b []byte) (int, error) {
	return sw.writer.Write(append([]byte(termite.TermControlEraseLine), b...))
}

// configureSpinner starts a spinner progress indicator on Stdout and hides the cursor.
// Reconfigures logrus globally with a formatter that prevents interferences.
// Call the returned cancel function to reverse the configuration
func configureSpinner() context.CancelFunc {
	cursor := termite.NewCursor(StdoutWriter)
	cursor.Hide()

	spinner := termite.NewSpinner(StdoutWriter, "", time.Millisecond*100)
	cancel, _ := spinner.Start()

	return func() {
		cursor.Show()
		cancel()
	}
}

func configureOutput(cmd *cobra.Command) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)
	var level = log.InfoLevel
	var writer = StdoutWriter

	if silent && debug {
		CheckUserArgFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		level = log.PanicLevel
		writer = StderrWriter
	}
	if debug {
		level = log.DebugLevel
		writer = StdoutWriter
	}

	log.StandardLogger().SetLevel(level)
	log.StandardLogger().SetOutput(writer)
}

func configureNonInteractiveOutput(cmd *cobra.Command) (cancel context.CancelFunc) {
	cancel = func() {}

	configureOutput(cmd)

	if termite.Tty && !GetBool(cmd, ArgNamePipeStdout) {
		cancelSpinner := configureSpinner()
		textFormatter := &log.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      true, // TTY mode
		}
		log.SetFormatter(textFormatter)
		log.StandardLogger().SetOutput(&alwaysRewritingWriter{log.StandardLogger().Out})

		cancel = cancelSpinner
	}

	return cancel
}
