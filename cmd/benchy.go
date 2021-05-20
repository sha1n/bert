package main

import (
	"fmt"

	"github.com/sha1n/benchy/internal/cli"
	"github.com/spf13/cobra"
)

var (
	// ProgramName : passed from build environment
	ProgramName string
	// Build : passed from build environment
	Build string
	// Version : passed from build environment
	Version string
)

func main() {
	var rootCmd = &cobra.Command{
		Use: ProgramName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, Version, Build),
		Example:      fmt.Sprintf("%s --%s <config file path>", ProgramName, cli.ArgNameConfig),
		SilenceUsage: false,
		Run:          cli.Run,
	}

	rootCmd.Flags().StringP(cli.ArgNameConfig, "c", "", `config file path. '~' will be expanded.`)

	// Reporting
	rootCmd.Flags().StringP(cli.ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)
	rootCmd.Flags().StringP(cli.ArgNameFormat, "f", "txt", `summary format. One of: 'txt', 'csv', 'csv/raw'
txt 		- plain text. designed to be used in your terminal
csv 		- CSV in which each row represents a scenario and contians calculated stats for that scenario
csv/raw	- CSV in which each row represents a raw trace event. useful if you want to import to a spreadsheet for further analysis`,
	)
	rootCmd.Flags().StringSliceP(cli.ArgNameLabel, "l", []string{}, `labels to attach to be included in the benchmark report`)
	rootCmd.Flags().BoolP(cli.ArgNameHeaders, "", true, `in tabular formats, whether to include headers in the report`)

	// Stdout
	rootCmd.Flags().BoolP(cli.ArgNamePipeStdout, "", true, `redirects external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP(cli.ArgNamePipeStderr, "", true, `redirects external commands standard error to benchy's standard error`)
	rootCmd.Flags().BoolP(cli.ArgNameDebug, "d", false, `logs extra debug information`)
	rootCmd.Flags().BoolP(cli.ArgNameSilent, "s", false, `logs only fatal errors`)

	cobra.MarkFlagRequired(rootCmd.Flags(), cli.ArgNameConfig)
	cobra.MarkFlagFilename(rootCmd.Flags(), cli.ArgNameConfig)
	cobra.MarkFlagFilename(rootCmd.Flags(), cli.ArgNameOutputFile)

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	_ = rootCmd.Execute()
}
