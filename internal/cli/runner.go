package cli

import (
	"bufio"
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
	var outFile *os.File
	var spec *api.BenchmarkSpec

	log.Info("Starting benchy...")

	if debug, _ := cmd.Flags().GetBool(ArgNameDebug); debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}

	if outFile, err = resolveOutputFile(cmd); err == nil {
		specFilePath, _ := cmd.Flags().GetString(ArgNameConfig)

		if spec, err = pkg.LoadSpec(specFilePath); err == nil {
			execCtx := resolveExecutionContext(cmd, spec)
			reportCtx := resolveReportContext(cmd)

			reportHandler := resolveReportHandler(cmd, spec, reportCtx, outFile)
			reportHandler.Subscribe(execCtx.Tracer.Stream())

			pkg.Execute(spec, execCtx)

			err = reportHandler.Finalize()
		}
	}

	checkFatal(err)
}

func checkFatal(err error) {
	if err != nil {
		log.Errorf("Failed to execute benchark. Error: %s", err.Error())
		log.Exit(1)
	}
}

func resolveReportHandler(cmd *cobra.Command, spec *api.BenchmarkSpec, reportCtx *api.ReportContext, outFile *os.File) api.ReportHandler {
	writer := bufio.NewWriter(outFile)
	writeReportFn := resolveReportWriter(cmd, writer)

	return pkg.NewSummaryReportHandler(spec, reportCtx, writeReportFn)
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

func resolveReportWriter(cmd *cobra.Command, writer *bufio.Writer) api.WriteReportFn {
	resolvedWriterFn := func() api.WriteReportFn {
		if reportFormat, _ := cmd.Flags().GetString(ArgNameFormat); reportFormat == "csv" {
			return internal.NewCsvReportWriter(writer)
		}

		var colorsOn = false
		if file, _ := cmd.Flags().GetString(ArgNameOutputFile); file == "" {
			colorsOn = true
		}

		return internal.NewTextReportWriter(writer, colorsOn)
	}()

	return internal.WriteReportFnFor(resolvedWriterFn)
}

func resolveOutputFile(cmd *cobra.Command) (outputFile *os.File, err error) {
	outputFile = os.Stdout
	if outputFilePath, _ := cmd.Flags().GetString(ArgNameOutputFile); outputFilePath != "" {
		resolvedfilePath := expandPath(outputFilePath)
		return os.OpenFile(resolvedfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	return outputFile, nil
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
