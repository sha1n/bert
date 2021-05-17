package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/internal"
	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	// CLIArgConfig : program arg name
	CLIArgConfig = "config"
	// CLIArgOutputFile : program arg name
	CLIArgOutputFile = "out-file"
	// CLIArgFormat : program arg name
	CLIArgFormat = "format"
	// CLIArgPipeStdout : program arg name
	CLIArgPipeStdout = "pipe-stdout"
	// CLIArgPipeStderr : program arg name
	CLIArgPipeStderr = "pipe-stderr"
	// CLIArgDebug : program arg name
	CLIArgDebug = "debug"
)

var (
	// ProgramName : passed from build environment
	ProgramName string
	// Build : passed from build environment
	Build string
	// Version : passed from build environment
	Version string
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var rootCmd = &cobra.Command{
		Use: ProgramName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, Version, Build),
		Example:      fmt.Sprintf("%s --%s <config file path>", ProgramName, CLIArgConfig),
		SilenceUsage: false,
		Run:          doRun,
	}

	rootCmd.Flags().StringP(CLIArgConfig, "c", "", `config file path`)
	rootCmd.Flags().StringP(CLIArgOutputFile, "o", "", `output file path`)
	rootCmd.Flags().StringP(CLIArgFormat, "f", "txt", `summary format. One of: 'txt', 'csv' (default: txt)`)
	rootCmd.Flags().BoolP(CLIArgPipeStdout, "", true, `redirects external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP(CLIArgPipeStderr, "", true, `redirects external commands standard error to benchy's standard error`)
	rootCmd.Flags().BoolP(CLIArgDebug, "d", false, `logs extra debug information`)

	cobra.MarkFlagRequired(rootCmd.Flags(), "config")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	_ = rootCmd.Execute()
}

func doRun(cmd *cobra.Command, args []string) {
	var err error
	var outputFile *os.File

	if debug, _ := cmd.Flags().GetBool(CLIArgDebug); debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}

	if outputFile, err = resolveOutputFile(cmd); err == nil {
		writeReportFn := resolveReportWriter(cmd, outputFile)
		ctx := resolveExecutionContext(cmd)
		specFilePath, _ := cmd.Flags().GetString(CLIArgConfig)
		if err = pkg.Run(specFilePath, ctx, writeReportFn); err != nil {
			log.Errorf("Failed to execute benchark. Error: %s", err.Error())
		}
	}
}

func resolveExecutionContext(cmd *cobra.Command) *api.ExecutionContext {
	pipeStdOut, _ := cmd.Flags().GetBool(CLIArgPipeStdout)
	pipeStdErr, _ := cmd.Flags().GetBool(CLIArgPipeStderr)

	return api.NewExecutionContext(pkg.NewTracer(), pkg.NewCommandExecutor(pipeStdOut, pipeStdErr))
}

func resolveReportWriter(cmd *cobra.Command, outputFile *os.File) api.WriteReportFn {
	resolvedWriterFn := func() api.WriteReportFn {
		writer := bufio.NewWriter(outputFile)
		if outputFilePath, _ := cmd.Flags().GetString(CLIArgFormat); outputFilePath == "csv" {
			return internal.NewCsvReportWriter(writer)
		}

		var colorsOn = false
		if file, _ := cmd.Flags().GetString(CLIArgOutputFile); file == "" {
			colorsOn = true
		}

		return internal.NewTextReportWriter(writer, colorsOn)
	}()

	return internal.WriteReportFnFor(resolvedWriterFn)
}

func resolveOutputFile(cmd *cobra.Command) (outputFile *os.File, err error) {
	outputFile = os.Stdout
	if outputFilePath, _ := cmd.Flags().GetString(CLIArgOutputFile); outputFilePath != "" {
		return os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	return outputFile, nil
}
