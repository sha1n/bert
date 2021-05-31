package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/internal/report"
	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the main command parse
func NewRootCommand(programName, version, build string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: programName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, version, build),
		Example:      fmt.Sprintf("%s --%s <config file path>", programName, ArgNameConfig),
		SilenceUsage: false,
		Run:          Run,
	}

	rootCmd.Flags().StringP(ArgNameConfig, "c", "", `config file path. '~' will be expanded.`)

	// Reporting
	rootCmd.Flags().StringP(ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)
	rootCmd.Flags().StringP(ArgNameFormat, "f", "txt", `summary format. One of: 'txt', 'md', 'md/raw', 'csv', 'csv/raw'
txt     - plain text. designed to be used in your terminal
md      - markdown table. similar to CSV but writes in markdown table format
md/raw  - markdown table in which each row represents a raw trace event.
csv     - CSV in which each row represents a scenario and contains calculated stats for that scenario
csv/raw - CSV in which each row represents a raw trace event. useful if you want to import to a spreadsheet for further analysis`,
	)
	rootCmd.Flags().StringSliceP(ArgNameLabel, "l", []string{}, `labels to attach to be included in the benchmark report`)
	rootCmd.Flags().BoolP(ArgNameHeaders, "", true, `in tabular formats, whether to include headers in the report`)

	// Stdout
	rootCmd.Flags().BoolP(ArgNamePipeStdout, "", false, `pipes external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP(ArgNamePipeStderr, "", false, `pipes external commands standard error to benchy's standard error`)

	rootCmd.PersistentFlags().BoolP(ArgNameDebug, "d", false, `logs extra debug information`)
	rootCmd.PersistentFlags().BoolP(ArgNameSilent, "s", false, `logs only fatal errors`)
	
	rootCmd.PersistentFlags().StringSliceP(ArgNameExperimental, "", []string{}, `enables a named experimental feature`)

	_ = rootCmd.MarkFlagRequired(ArgNameConfig)
	_ = rootCmd.MarkFlagFilename(ArgNameConfig, "yml", "yaml", "json")
	_ = rootCmd.MarkFlagFilename(ArgNameOutputFile, "txt", "csv", "md")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	return rootCmd
}

// Run parses CLI arguments and runs the benchmark process
func Run(cmd *cobra.Command, args []string) {
	var err error
	var closer io.Closer
	defer configureNonInteractiveOutput(cmd)()

	log.Info("Starting benchy...")

	var spec *api.BenchmarkSpec
	spec, err = loadSpec(cmd)
	CheckBenchmarkInitFatal(err)

	var reportHandler api.ReportHandler
	reportHandler, closer, err = resolveReportHandler(cmd, spec)
	defer closer.Close()

	if err == nil {
		tracer := pkg.NewTracer(spec.Executions * len(spec.Scenarios))
		reportHandler.Subscribe(tracer.Stream())

		pkg.Execute(spec, resolveExecutionContext(cmd, tracer))

		err = reportHandler.Finalize()
	}

	CheckFatal(err)
}

func loadSpec(cmd *cobra.Command) (spec *api.BenchmarkSpec, err error) {
	var filePath string
	filePath = GetString(cmd, ArgNameConfig)
	filePath, err = filepath.Abs(expandPath(filePath))

	if err == nil {
		_, err = os.Stat(expandPath(filePath))
		exists := !os.IsNotExist(err)

		if err == nil && exists {
			return pkg.LoadSpec(filePath)
		}

		err = fmt.Errorf("the file '%s' does not exist, or is not accessible", filePath)
	}

	return spec, err
}

func resolveReportHandler(cmd *cobra.Command, spec *api.BenchmarkSpec) (handler api.ReportHandler, closer io.Closer, err error) {
	reportCtx := resolveReportContext(cmd)
	writeCloser := ResolveOutputArg(cmd, ArgNameOutputFile)
	writer := bufio.NewWriter(writeCloser)

	switch reportFormat := GetString(cmd, ArgNameFormat); reportFormat {
	case ArgValueReportFormatMarkdownRaw:
		streamReportWriter := report.NewMarkdownStreamReportWriter(writer, reportCtx)
		handler = pkg.NewStreamReportHandler(spec, reportCtx, streamReportWriter.Handle)

	case ArgValueReportFormatCsvRaw:
		streamReportWriter := report.NewCsvStreamReportWriter(writer, reportCtx)
		handler = pkg.NewStreamReportHandler(spec, reportCtx, streamReportWriter.Handle)

	case ArgValueReportFormatMarkdown:
		handler = pkg.NewSummaryReportHandler(spec, reportCtx, report.NewMarkdownSummaryReportWriter(writer))

	case ArgValueReportFormatCsv:
		handler = pkg.NewSummaryReportHandler(spec, reportCtx, report.NewCsvReportWriter(writer))

	case ArgValueReportFormatTxt:
		var colorsOn = false
		if GetString(cmd, ArgNameOutputFile) == "" {
			colorsOn = true
		}

		handler = pkg.NewSummaryReportHandler(spec, reportCtx, report.NewTextReportWriter(writer, colorsOn))

	default:
		err = fmt.Errorf("Invalid report format '%s'", reportFormat)
	}

	return handler, writeCloser, err
}

func resolveReportContext(cmd *cobra.Command) *api.ReportContext {
	return &api.ReportContext{
		Labels:         GetStringSlice(cmd, ArgNameLabel),
		IncludeHeaders: GetBool(cmd, ArgNameHeaders),
	}
}

func resolveExecutionContext(cmd *cobra.Command, tracer api.Tracer) *api.ExecutionContext {
	pipeStdOut := GetBool(cmd, ArgNamePipeStdout)
	pipeStdErr := GetBool(cmd, ArgNamePipeStderr)

	return api.NewExecutionContext(tracer, pkg.NewCommandExecutor(pipeStdOut, pipeStdErr))
}
