package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const (
	// ArgNameConfig : program arg name
	ArgNameConfig = "config"
	// ArgNameOutputFile : program arg name
	ArgNameOutputFile = "out-file"
	// ArgNameFormat : program arg name
	ArgNameFormat = "format"
	// ArgNamePipeStdout : program arg name
	ArgNamePipeStdout = "pipe-stdout"
	// ArgNamePipeStderr : program arg name
	ArgNamePipeStderr = "pipe-stderr"
	// ArgNameDebug : program arg name
	ArgNameDebug = "debug"
	// ArgNameSilent : program arg name
	ArgNameSilent = "silent"
	// ArgNameLabel : program arg name
	ArgNameLabel = "label"
	// ArgNameHeaders : program arg name
	ArgNameHeaders = "headers"

	// ArgValueReportFormatTxt : Plain text report format arg value
	ArgValueReportFormatTxt = "txt"
	// ArgValueReportFormatCsv : CSV report format arg value
	ArgValueReportFormatCsv = "csv"
	// ArgValueReportFormatCsvRaw : CSV raw data report format value
	ArgValueReportFormatCsvRaw = "csv/raw"
	// ArgValueReportFormatMarkdown : Markdown report format arg value
	ArgValueReportFormatMarkdown = "md"
	// ArgValueReportFormatMarkdownRaw : Markdown report format arg value
	ArgValueReportFormatMarkdownRaw = "md/raw"
)

// ResolveOutputArg resolves an output file argument based on user input.
// If the specified argument is empty, stdout is returned.
func ResolveOutputArg(cmd *cobra.Command, name string) io.WriteCloser {
	var outputFile io.WriteCloser = stdOutNonClosingWriteCloser{}
	var err error = nil

	if outputFilePath := GetString(cmd, name); outputFilePath != "" {
		resolvedfilePath := expandPath(outputFilePath)
		outputFile, err = os.OpenFile(resolvedfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	CheckBenchmarkInitFatal(err)

	return outputFile
}

// GetString tries to get a user argument. Handles errors as fatal.
func GetString(cmd *cobra.Command, name string) string {
	v, err := cmd.Flags().GetString(name)
	CheckUserArgFatal(err)

	return v
}

// GetBool tries to get a user argument. Handles errors as fatal.
func GetBool(cmd *cobra.Command, name string) bool {
	var v bool
	var err error
	if v, err = cmd.Flags().GetBool(name); err != nil {
		v, err = cmd.PersistentFlags().GetBool(name)
	}
	CheckUserArgFatal(err)

	return v
}

// GetStringSlice tries to get a user argument. Handles errors as fatal.
func GetStringSlice(cmd *cobra.Command, name string) []string {
	v, err := cmd.Flags().GetStringSlice(ArgNameLabel)
	CheckUserArgFatal(err)

	return v
}

// TODO this has been copied from pgk/command_exec.go. Maybe share or use an existing implementation if exists.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if p, err := os.UserHomeDir(); err == nil {
			return filepath.Join(p, path[1:])
		}
		log.Warnf("Failed to resolve user home for path '%s'", path)
	}

	return path
}

// stdOutNonClosingWriteCloser a wrapper around os.Stdout that implements the io.WriteCloser interface but never closes the file
type stdOutNonClosingWriteCloser struct{}

// Write forwards the call to standard output
func (wc stdOutNonClosingWriteCloser) Write(b []byte) (int, error) {
	logger := log.StandardLogger()
	return logger.Out.Write(b)
}

// Close NOOP
func (wc stdOutNonClosingWriteCloser) Close() error {
	return nil
}
