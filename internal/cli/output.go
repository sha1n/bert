package cli

import (
	"context"
	"errors"
	"io"
	"syscall"
	"time"

	"github.com/fatih/color"
	clibos "github.com/sha1n/clib/pkg/os"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const richOutputExperimentName = "rich_output"

var (
	printfRed   = color.New(color.FgRed).Printf
	printRed    = color.New(color.FgRed).Print
	sprintRed   = color.New(color.FgRed).Sprint
	sprintGreen = color.New(color.FgGreen).Sprint
	sprintBold  = color.New(color.Bold).Sprint
)

func configureDefaultOutput(cmd *cobra.Command, ctx IOContext) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)
	var level = log.InfoLevel

	if silent && debug {
		CheckUserArgFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		level = log.PanicLevel
	}
	if debug {
		level = log.DebugLevel
	}
	if ctx.Tty {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      true,
		})
	}

	log.StandardLogger().SetLevel(level)
	log.StandardLogger().SetOutput(ctx.StderrWriter)
}

// configureSpinner starts a spinner progress indicator on Stdout and hides the cursor.
// Reconfigures logrus globally with a formatter that prevents interferences.
// Call the returned cancel function to reverse the configuration
func configureSpinner(writer io.Writer) context.CancelFunc {
	cursor := termite.NewCursor(writer)
	cursor.Hide()

	spinner := termite.NewSpinner(writer, "", time.Millisecond*100, termite.DefaultSpinnerFormatter())
	cancel, _ := spinner.Start()

	onShutdownSignal := func() {
		log.Debugf("Received OS shutdown signal")
		cursor.Show()
		log.Debugf("Exiting!")
		log.Exit(1)
	}

	clibos.RegisterShutdownHook(clibos.NewSignalHook(syscall.SIGINT, onShutdownSignal))
	clibos.RegisterShutdownHook(clibos.NewSignalHook(syscall.SIGTERM, onShutdownSignal))

	return func() {
		cursor.Show()
		cancel()
	}
}

func configureRichOutput(cmd *cobra.Command, ctx IOContext) (cancel context.CancelFunc) {
	cancel = func() {}

	configureDefaultOutput(cmd, ctx)

	if ctx.Tty {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      true, // TTY mode
		})
	}

	if IsExperimentEnabled(cmd, richOutputExperimentName) {
		cancel = configureSpinner(ctx.StdoutWriter)

		log.StandardLogger().SetOutput(&alwaysRewritingWriter{log.StandardLogger().Out})
	}

	return cancel
}

// this writer should be used when background threads write messages without a new line suffix to stdout.
// this is common with progress indicartors from the termite package.
type alwaysRewritingWriter struct {
	writer io.Writer
}

func (sw *alwaysRewritingWriter) Write(b []byte) (int, error) {
	return sw.writer.Write(append([]byte(termite.TermControlEraseLine), b...))
}
