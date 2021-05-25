package cli

import (
	"io"
	"os"
)

// StdinReader is used as a replacement for the global os.StdinReader, primarily to make testing easier
var StdinReader io.Reader = os.Stdin
