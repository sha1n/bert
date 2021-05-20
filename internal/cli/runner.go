package cli

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/internal/report"
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
	var outputFile *os.File
	var spec *api.BenchmarkSpec

	log.Info("Starting benchy...")

	if debug, _ := cmd.Flags().GetBool(ArgNameDebug); debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}

	if outputFile, err = resolveOutputFile(cmd); err == nil {
		specFilePath, _ := cmd.Flags().GetString(ArgNameConfig)

		if spec, err = pkg.LoadSpec(specFilePath); err == nil {
			ctx := resolveExecutionContext(cmd)

			summary := pkg.Execute(spec, ctx)

			writeReportFn := resolveReportWriter(cmd, outputFile)
			reportCtx := resolveReportContext(cmd)

			writeReportFn(summary, spec, reportCtx)

		}
	}

	checkFatal(err)
}

func checkFatal(err error) {
	log.Errorf("Failed to execute benchark. Error: %s", err.Error())
	log.Exit(1)
}

func resolveReportContext(cmd *cobra.Command) *api.ReportContext {
	labels, _ := cmd.Flags().GetStringSlice(ArgNameLabel)
	includeHeaders, _ := cmd.Flags().GetBool(ArgNameHeaders)

	return &api.ReportContext{
		Labels:         labels,
		IncludeHeaders: includeHeaders,
	}
}

func resolveExecutionContext(cmd *cobra.Command) *api.ExecutionContext {
	pipeStdOut, _ := cmd.Flags().GetBool(ArgNamePipeStdout)
	pipeStdErr, _ := cmd.Flags().GetBool(ArgNamePipeStderr)

	return api.NewExecutionContext(
		pkg.NewTracer(),
		pkg.NewCommandExecutor(pipeStdOut, pipeStdErr),
	)
}

func resolveReportWriter(cmd *cobra.Command, outputFile *os.File) api.WriteReportFn {
	resolvedWriterFn := func() api.WriteReportFn {
		writer := bufio.NewWriter(outputFile)
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
