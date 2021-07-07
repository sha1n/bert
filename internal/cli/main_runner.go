package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/internal/report"
	"github.com/sha1n/bert/pkg/exec"
	"github.com/sha1n/bert/pkg/osutil"
	"github.com/sha1n/bert/pkg/reporthandlers"
	"github.com/sha1n/bert/pkg/specs"
	"github.com/sha1n/bert/pkg/ui"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const bert = `

                      \WWW/
                      /   \
                     /wwwww\
                   _|  o_o  |_
                  (_   / \   _)
                    |  \_/  |
                    : ~~~~~ :
                     \_____/
                     [     ]
                      """""
               ____            _   
              |  _ \          | |  
              | |_) | ___ _ __| |_ 
              |  _ < / _ \ '__| __|
              | |_) |  __/ |  | |_ 
              |____/ \___|_|   \__|

`

// NewRootCommand creates the main command parse
func NewRootCommand(programName, version, build string, ctx api.IOContext) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: programName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, version, build),
		Example: fmt.Sprintf(`
	%s --%s <config file path>
	%s 'command -optA' 'command -optB' --%s 100`, programName, ArgNameConfig, programName, ArgNameExecutions),

		SilenceUsage: false,
		Args:         validatePositionalArgs,
		Run:          runFn(ctx),
	}

	rootCmd.Flags().StringP(ArgNameConfig, "c", "", `config file path. '~' will be expanded.`)
	rootCmd.Flags().IntP(ArgNameExecutions, "e", 0, `the number of executions per scenario.
required when no configuration file is provided. 
when specified with a configuration file, this argument has priority.`)

	// Reporting
	rootCmd.Flags().StringP(ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)
	rootCmd.Flags().StringP(ArgNameFormat, "f", "txt", `summary format. One of: 'txt', 'md', 'md/raw', 'csv', 'csv/raw'
txt     - plain text. designed to be used in your terminal.
json    - JSON document. each object represents a scenario and contains calculated stats for that scenario.
csv     - CSV document. each row represents a scenario and contains calculated stats for that scenario.
csv/raw - CSV document in which each row represents a raw trace event. useful if you want to import to a spreadsheet for further analysis.
md      - markdown table. similar to CSV but writes in markdown table format.
md/raw  - markdown table in which each row represents a raw trace event.`,
	)
	rootCmd.Flags().StringSliceP(ArgNameLabel, "l", []string{}, `labels to attach to be included in the benchmark report.`)
	rootCmd.Flags().Bool(ArgNameHeaders, true, `in tabular formats, whether to include headers in the report.`)
	rootCmd.Flags().Bool(ArgReportUTCDate, false, `whether to use UTC date.`)

	// Stdout
	rootCmd.Flags().Bool(ArgNamePipeStdout, false, `pipes external commands standard out to bert's standard out.`)
	rootCmd.Flags().Bool(ArgNamePipeStderr, false, `pipes external commands standard error to bert's standard error.`)

	rootCmd.PersistentFlags().BoolP(ArgNameDebug, "d", false, `runs the program in debug mode.`)
	rootCmd.PersistentFlags().BoolP(ArgNameSilent, "s", false, `logs only fatal errors.`)

	rootCmd.PersistentFlags().StringSlice(ArgNameExperimental, []string{}, `enables a named experimental features.`)

	_ = rootCmd.MarkFlagFilename(ArgNameConfig, "yml", "yaml", "json")
	_ = rootCmd.MarkFlagFilename(ArgNameOutputFile, "txt", "csv", "md", "json")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)
	rootCmd.SetUsageTemplate(rootCmd.UsageTemplate() + bert)

	return rootCmd
}

// runFn returns a function that parses CLI arguments and runs the benchmark process with the specified IOContext
func runFn(ctx api.IOContext) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		var err error
		var closer io.Closer
		configureOutput(cmd, log.ErrorLevel, ctx)

		log.Info("Starting bert...")

		var spec api.BenchmarkSpec
		spec, err = loadSpec(cmd, args)
		CheckBenchmarkInitFatal(err)

		var reportHandler api.ReportHandler
		reportHandler, closer, err = resolveReportHandler(cmd, spec, ctx)
		defer closer.Close()

		if err == nil {
			tracer := exec.NewTracer(spec.Executions * len(spec.Scenarios))
			reportHandler.Subscribe(tracer.Stream())

			log.Info("Executing...")
			exec.Execute(spec, resolveExecutionContext(cmd, spec, ctx, tracer))

			log.Info("Finalizing report...")
			err = reportHandler.Finalize()

			log.Info("Done")
		}

		CheckFatal(err)
	}

}

func loadSpec(cmd *cobra.Command, args []string) (spec api.BenchmarkSpec, err error) {
	executions := GetInt(cmd, ArgNameExecutions)

	if len(args) > 0 { // positional args are used for ad-hoc config
		commands := []api.CommandSpec{}
		for i := range args {
			commands = append(commands, api.CommandSpec{
				Cmd: parseCommand(strings.Trim(args[i], "'\"")),
			})
		}
		spec, err = specs.CreateSpecFrom(executions, false, commands...)

	} else {
		var filePath string
		filePath = GetString(cmd, ArgNameConfig)
		filePath, err = filepath.Abs(osutil.ExpandUserPath(filePath))

		if err == nil {
			_, err = os.Stat(osutil.ExpandUserPath(filePath))
			exists := !os.IsNotExist(err)

			if err != nil || !exists {
				err = fmt.Errorf("the file '%s' does not exist, or is not accessible", filePath)
				return
			}

			if spec, err = specs.LoadSpec(filePath); err != nil {
				return
			}
		}
	}

	// Override executions if specified
	if executions > 0 {
		spec.Executions = executions
	}

	return spec, err
}

func resolveReportHandler(cmd *cobra.Command, spec api.BenchmarkSpec, ctx api.IOContext) (handler api.ReportHandler, closer io.Closer, err error) {
	reportCtx := resolveReportContext(cmd)
	writeCloser := ResolveOutputArg(cmd, ArgNameOutputFile, ctx)
	writer := writeCloser

	switch reportFormat := GetString(cmd, ArgNameFormat); reportFormat {
	case ArgValueReportFormatMarkdownRaw:
		streamReportWriter := report.NewMarkdownStreamReportWriter(writer, reportCtx)
		handler = report_handlers.NewStreamReportHandler(spec, reportCtx, streamReportWriter.Handle)

	case ArgValueReportFormatCsvRaw:
		streamReportWriter := report.NewCsvStreamReportWriter(writer, reportCtx)
		handler = report_handlers.NewStreamReportHandler(spec, reportCtx, streamReportWriter.Handle)

	case ArgValueReportFormatMarkdown:
		handler = report_handlers.NewSummaryReportHandler(spec, reportCtx, report.NewMarkdownSummaryReportWriter(writer))

	case ArgValueReportFormatCsv:
		handler = report_handlers.NewSummaryReportHandler(spec, reportCtx, report.NewCsvReportWriter(writer))

	case ArgValueReportFormatJSON:
		handler = report_handlers.NewSummaryReportHandler(spec, reportCtx, report.NewJSONReportWriter(writer))

	case ArgValueReportFormatTxt:
		var colorsOn = false
		if GetString(cmd, ArgNameOutputFile) == "" {
			colorsOn = true
		}

		handler = report_handlers.NewSummaryReportHandler(spec, reportCtx, report.NewTextReportWriter(writer, colorsOn))

	default:
		err = fmt.Errorf("Invalid report format '%s'", reportFormat)
	}

	return handler, writeCloser, err
}

func resolveReportContext(cmd *cobra.Command) api.ReportContext {
	return api.ReportContext{
		Labels:         GetStringSlice(cmd, ArgNameLabel),
		IncludeHeaders: GetBool(cmd, ArgNameHeaders),
		UTCDate:        GetBool(cmd, ArgReportUTCDate),
	}
}

func resolveExecutionContext(cmd *cobra.Command, spec api.BenchmarkSpec, ctx api.IOContext, tracer api.Tracer) api.ExecutionContext {
	pipeStdOut := GetBool(cmd, ArgNamePipeStdout)
	pipeStdErr := GetBool(cmd, ArgNamePipeStderr)

	return api.NewExecutionContext(
		tracer,
		exec.NewCommandExecutor(pipeStdOut, pipeStdErr),
		resolveExecutionListener(cmd, spec, ctx),
	)
}

func resolveExecutionListener(cmd *cobra.Command, spec api.BenchmarkSpec, ctx api.IOContext) api.Listener {
	if enableTerminalGUI(cmd, ctx) {
		return ui.NewProgressView(spec, terminalDimensionsOrFake, ctx)
	}

	return ui.NewLoggingProgressListener()
}

func enableTerminalGUI(cmd *cobra.Command, ctx api.IOContext) bool {
	reportToFile := GetString(cmd, ArgNameOutputFile) != ""
	enableRichOut := reportToFile || !StreamingReportFormats[GetString(cmd, ArgNameFormat)]
	silentMode := GetBool(cmd, ArgNameSilent)
	debugMode := GetBool(cmd, ArgNameDebug)
	pipeOutputsMode := GetBool(cmd, ArgNamePipeStdout)
	pipeOutputsMode = pipeOutputsMode || GetBool(cmd, ArgNamePipeStderr)

	return ctx.Tty && enableRichOut && !(silentMode || debugMode || pipeOutputsMode)
}

func terminalDimensionsOrFake() (int, int) {
	if w, h, err := termite.GetTerminalDimensions(); err == nil {
		return w, h
	}
	return 0, 0

}

func validatePositionalArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		executions := GetInt(cmd, ArgNameExecutions)
		if executions < 1 {
			return fmt.Errorf("--%s is required", ArgNameExecutions)
		}
	} else {
		if outputFilePath := GetString(cmd, ArgNameConfig); outputFilePath == "" {
			return fmt.Errorf("either specify a configuration file with '--%s', or inline commands to benchmark", ArgNameConfig)
		}
	}
	return nil
}
