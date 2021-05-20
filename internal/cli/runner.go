package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sha1n/benchy/api"
	internal "github.com/sha1n/benchy/internal/report"
	"github.com/sha1n/benchy/pkg"
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
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Run parses CLI arguments and runs the benchmark process
func Run(cmd *cobra.Command, args []string) {
	var err error

	log.Info("Starting benchy...")

	handleDebug(cmd)

	spec := loadSpec(cmd)
	execCtx := resolveExecutionContext(cmd, spec)
	if reportHandler, err := resolveReportHandler(cmd, spec); err == nil {
		reportHandler.Subscribe(execCtx.Tracer.Stream())

		pkg.Execute(spec, execCtx)

		err = reportHandler.Finalize()
	}

	checkFatal(err)
}

func checkFatal(err error) {
	if err != nil {
		log.Errorf("Failed to execute benchark. Error: %s", err.Error())
		log.Exit(1)
	}
}

func handleDebug(cmd *cobra.Command) {
	if debug, _ := cmd.Flags().GetBool(ArgNameDebug); debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}
}

func loadSpec(cmd *cobra.Command) *api.BenchmarkSpec {
	specFilePath, _ := cmd.Flags().GetString(ArgNameConfig)
	spec, err := pkg.LoadSpec(specFilePath)
	checkFatal(err)

	return spec
}

func resolveReportHandler(cmd *cobra.Command, spec *api.BenchmarkSpec) (handler api.ReportHandler, err error) {
	reportCtx := resolveReportContext(cmd)
	outFile := resolveOutputFile(cmd)
	writer := bufio.NewWriter(outFile)

	switch reportFormat, _ := cmd.Flags().GetString(ArgNameFormat); reportFormat {
	case ArgValueReportFormatCsvRaw:
		streamReportWriter := internal.NewCsvStreamReportWriter(writer, reportCtx)
		handler = pkg.NewStreamReportHandler(spec, reportCtx, streamReportWriter.Handle)

	case ArgValueReportFormatCsv:
		handler = pkg.NewSummaryReportHandler(spec, reportCtx, internal.NewCsvReportWriter(writer))

	case ArgValueReportFormatTxt:
		var colorsOn = false
		if file, _ := cmd.Flags().GetString(ArgNameOutputFile); file == "" {
			colorsOn = true
		}

		handler = pkg.NewSummaryReportHandler(spec, reportCtx, internal.NewTextReportWriter(writer, colorsOn))

	default:
		err = fmt.Errorf("Invalid report format '%s'", reportFormat)
	}

	return handler, err
}

func resolveReportContext(cmd *cobra.Command) *api.ReportContext {
	labels, _ := cmd.Flags().GetStringSlice(ArgNameLabel)
	includeHeaders, _ := cmd.Flags().GetBool(ArgNameHeaders)

	return &api.ReportContext{
		Labels:         labels,
		IncludeHeaders: includeHeaders,
	}
}

func resolveExecutionContext(cmd *cobra.Command, spec *api.BenchmarkSpec) *api.ExecutionContext {
	pipeStdOut, _ := cmd.Flags().GetBool(ArgNamePipeStdout)
	pipeStdErr, _ := cmd.Flags().GetBool(ArgNamePipeStderr)

	return api.NewExecutionContext(
		pkg.NewTracer(spec.Executions*len(spec.Scenarios)),
		pkg.NewCommandExecutor(pipeStdOut, pipeStdErr),
	)
}

func resolveOutputFile(cmd *cobra.Command) *os.File {
	var outputFile = os.Stdout
	var err error = nil

	if outputFilePath, _ := cmd.Flags().GetString(ArgNameOutputFile); outputFilePath != "" {
		resolvedfilePath := expandPath(outputFilePath)
		outputFile, err = os.OpenFile(resolvedfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	checkFatal(err)

	return outputFile
}

// FIXME this has been copied from pgk/command_exec.go. Maybe share or use an existing implementation if exists.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if p, err := os.UserHomeDir(); err == nil {
			return filepath.Join(p, path[1:])
		}
		log.Warnf("Failed to resolve user home for path '%s'", path)
	}

	return path
}
