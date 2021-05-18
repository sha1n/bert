package cli

import (
	"bufio"
	"os"

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
			labels, _ := cmd.Flags().GetStringSlice(ArgNameLabel)
			writeReportFn(summary, spec, &api.ReportContext{Labels: labels})
		} else {
			log.Errorf("Failed to execute benchark. Error: %s", err.Error())
		}
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
		if outputFilePath, _ := cmd.Flags().GetString(ArgNameFormat); outputFilePath == "csv" {
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
		return os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	return outputFile, nil
}
