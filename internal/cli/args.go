package cli

import (
	"io"
	"os"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg"

	"github.com/spf13/cobra"
)

const (
	// ArgNameConfig : program arg name
	ArgNameConfig = "config"
	// ArgNameExecutions : program arg name
	ArgNameExecutions = "executions"
	// ArgNameOutputFile : program arg name
	ArgNameOutputFile = "out-file"
	// ArgNameConfigExample : program arg name
	ArgNameConfigExample = "example"
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

	// ArgNameExperimental : enables a named experimental feature
	ArgNameExperimental = "experimental"

	// ArgNameLabel : program arg name
	ArgNameLabel = "label"
	// ArgNameHeaders : program arg name
	ArgNameHeaders = "headers"

	// ArgReportUTCDate : specifies that reports should report UTC time
	ArgReportUTCDate = "utc-date"
	// ArgValueReportFormatTxt : Plain text report format arg value
	ArgValueReportFormatTxt = "txt"
	// ArgValueReportFormatCsv : CSV report format arg value
	ArgValueReportFormatCsv = "csv"
	// ArgValueReportFormatJSON : JSON report format arg value
	ArgValueReportFormatJSON = "json"
	// ArgValueReportFormatCsvRaw : CSV raw data report format value
	ArgValueReportFormatCsvRaw = "csv/raw"
	// ArgValueReportFormatMarkdown : Markdown report format arg value
	ArgValueReportFormatMarkdown = "md"
	// ArgValueReportFormatMarkdownRaw : Markdown report format arg value
	ArgValueReportFormatMarkdownRaw = "md/raw"
)

// StreamingReportFormats a slice containing the values that represent report formats that are reporting in streaming
var StreamingReportFormats = map[string]bool{ArgValueReportFormatCsvRaw: true, ArgValueReportFormatMarkdownRaw: true}

// ResolveOutputArg resolves an output file argument based on user input.
// If the specified argument is empty, stdout is returned.
func ResolveOutputArg(cmd *cobra.Command, name string, ctx api.IOContext) io.WriteCloser {
	var outputFile io.WriteCloser = stdOutNonClosingWriteCloser{out: ctx.StdoutWriter}
	var err error = nil

	if outputFilePath := GetString(cmd, name); outputFilePath != "" {
		resolvedfilePath := pkg.ExpandUserPath(outputFilePath)
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

// GetInt tries to get a user argument. Handles errors as fatal.
func GetInt(cmd *cobra.Command, name string) int {
	v, err := cmd.Flags().GetInt(name)
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

// IsExperimentEnabled checks whether the specified experiment is enabled by the command line
func IsExperimentEnabled(cmd *cobra.Command, name string) bool {
	if slice, err := cmd.Flags().GetStringSlice(ArgNameExperimental); err == nil {
		for _, item := range slice {
			if item == name {
				return true
			}
		}
	}

	return false
}

// stdOutNonClosingWriteCloser a wrapper around os.Stdout that implements the io.WriteCloser interface but never closes the file
type stdOutNonClosingWriteCloser struct {
	out io.Writer
}

// Write forwards the call to standard output
func (wc stdOutNonClosingWriteCloser) Write(b []byte) (int, error) {
	return wc.out.Write(b)
}

// Close NOOP
func (wc stdOutNonClosingWriteCloser) Close() error {
	return nil
}
