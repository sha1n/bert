package cli

import (
	"io"
	"os"

	"github.com/sha1n/termite"
)

var (
	// StdinReader is used as a replacement for the global os.Stdin, primarily to make testing easier
	// Should not be used directly
	stdinReader io.Reader = os.Stdin

	// stdoutWriter should be used in place of Stdout
	// Should not be used directly
	stdoutWriter io.Writer = termite.StdoutWriter

	// stderrWriter should be used in place of Stderr
	// Should not be used directly
	stderrWriter io.Writer = termite.StderrWriter
)

// IOContext serves as a contextual accessor to key I/O elements.
// This enables more flexible design, easier and concurrent testing.
type IOContext struct {
	// StdoutWriter provides Stdout semantics
	StdoutWriter io.Writer
	// StderrWriter provides Stderr semantics
	StderrWriter io.Writer
	// StdinReader provides Stdin semantics
	StdinReader io.Reader
	// Tty whether or not if this process is connected to a terminal
	Tty bool
}

// NewIOContext returns a new IOContext populated with the global system I/O elements.
func NewIOContext() IOContext {
	return IOContext{
		StdoutWriter: stdoutWriter,
		StderrWriter: stderrWriter,
		StdinReader:  stdinReader,
		Tty:          termite.Tty,
	}
}
