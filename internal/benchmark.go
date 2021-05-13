package internal

import (
	"bufio"

	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// WriteReport writes a benchmark report.
type WriteReport = func(summary pkg.TracerSummary, config *BenchmarkSpec)

// Run accepts the program arguments and runs a benchmark.
func Run(cmd *cobra.Command, args []string) error {
	specFilePath, _ := cmd.Flags().GetString("config")
	pipeStdOut, _ := cmd.Flags().GetBool("pipe-stdout")
	pipeStdErr, _ := cmd.Flags().GetBool("pipe-stderr")

	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}

	executor := NewCommandExecutor(
		pipeStdOut,
		pipeStdErr,
	)

	return run(specFilePath, executor, writeReport)
}

func run(specFilePath string, executor CommandExecutor, writeReport WriteReport) (error error) {
	log.Info("Starting benchy...")

	spec, err := loadSpec(specFilePath)
	if err != nil {
		return err
	}

	summary := ExecuteBenchmark(spec, NewExecutionContext(pkg.NewTracer(), executor))
	writeReport(summary, spec)

	return error
}

func loadSpec(filePath string) (rtn *BenchmarkSpec, error error) {
	log.Infof("Loading benchmark specs from '%s'...", filePath)

	benchmark, err := Load(filePath)
	if err != nil {
		log.Error(err.Error())
		error = err
	}

	return benchmark, error
}

func writeReport(summary pkg.TracerSummary, config *BenchmarkSpec) {
	console := bufio.NewWriter(os.Stdout)
	writer := NewTextReportWriter(console)
	writer.Write(summary, config)
}
