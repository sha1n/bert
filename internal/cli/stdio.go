package cli

import (
	"io"
	"os"

	"github.com/sha1n/termite"
)

var (
	// StdinReader is used as a replacement for the global os.StdinReader, primarily to make testing easier
	StdinReader io.Reader = os.Stdin

	// StdoutWriter should be used in place of Stdout
	StdoutWriter io.Writer = termite.StdoutWriter

	// StderrWriter should be used in place of Stderr
	StderrWriter io.Writer = termite.StderrWriter
)
