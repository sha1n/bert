package internal

import (
	"bufio"
	"log"
	"os"

	"github.com/sha1n/benchy/pkg"
	"github.com/spf13/cobra"
)

type WriteReport = func(summary pkg.TracerSummary, config *BenchmarkSpec)

func Run(cmd *cobra.Command, args []string) error {
	specFilePath, _ := cmd.Flags().GetString("config")
	pipeStdOut, _ := cmd.Flags().GetBool("pipe-stdout")
	pipeStdErr, _ := cmd.Flags().GetBool("pipe-stderr")

	executor := NewCommandExecutor(
		pipeStdOut,
		pipeStdErr,
	)

	return run(specFilePath, executor, writeReport)
}

func run(specFilePath string, executor CommandExecutor, writeReport WriteReport) (error error) {
	log.Println("Starting benchy...")

	spec, err := loadSpec(specFilePath)
	if err != nil {
		return err
	}

	summary := ExecuteBenchmark(spec, NewExecutionContext(pkg.NewTracer(), executor))
	writeReport(summary, spec)

	return error
}

func loadSpec(filePath string) (rtn *BenchmarkSpec, error error) {
	log.Printf("Loading benchmark specs from '%s'...\r\n", filePath)

	benchmark, err := Load(filePath)
	if err != nil {
		log.Println(err)
		error = err
	}

	return benchmark, error
}

func writeReport(summary pkg.TracerSummary, config *BenchmarkSpec) {
	console := bufio.NewWriter(os.Stdout)
	writer := NewTextReportWriter(console)
	writer.Write(summary, config)
}
