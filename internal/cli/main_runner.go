package cli

import (
	"bufio"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/internal/report"
	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Run parses CLI arguments and runs the benchmark process
func Run(cmd *cobra.Command, args []string) {
	var err error
	configureLogger(cmd)

	log.Info("Starting benchy...")

	spec := loadSpec(cmd)
	var reportHandler api.ReportHandler
	if reportHandler, err = resolveReportHandler(cmd, spec); err == nil {
		tracer := pkg.NewTracer(spec.Executions * len(spec.Scenarios))
		reportHandler.Subscribe(tracer.Stream())

		pkg.Execute(spec, resolveExecutionContext(cmd, tracer))

		err = reportHandler.Finalize()
	}

	CheckBenchmarkInitFatal(err)
}

func configureLogger(cmd *cobra.Command) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)

	if silent && debug {
		CheckBenchmarkInitFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		log.StandardLogger().SetLevel(log.PanicLevel)
	}
	if debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}
}

func loadSpec(cmd *cobra.Command) (spec *api.BenchmarkSpec) {
	var filePath string
	var err error
	filePath = GetString(cmd, ArgNameConfig)
	if filePath, err = filepath.Abs(filePath); err == nil {
		spec, err = pkg.LoadSpec(filePath)
	}
	CheckBenchmarkInitFatal(err)

	return spec
}

func resolveReportHandler(cmd *cobra.Command, spec *api.BenchmarkSpec) (handler api.ReportHandler, err error) {
	reportCtx := resolveReportContext(cmd)
	outFile := ResolveOutputFileArg(cmd, ArgNameOutputFile)
	writer := bufio.NewWriter(outFile)

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

	return handler, err
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
