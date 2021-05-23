package main

import (
	"fmt"
	"os"

	"github.com/sha1n/benchy/cmd/subcmd"
	"github.com/sha1n/benchy/internal/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

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
	rootCmd.Flags().StringP(cli.ArgNameFormat, "f", "txt", `summary format. One of: 'txt', 'md', 'md/raw', 'csv', 'csv/raw'
txt     - plain text. designed to be used in your terminal
md      - markdown table. similar to CSV but writes in markdown table format
md/raw  - markdown table in which each row represents a raw trace event.
csv     - CSV in which each row represents a scenario and contains calculated stats for that scenario
csv/raw - CSV in which each row represents a raw trace event. useful if you want to import to a spreadsheet for further analysis`,
	)
	rootCmd.Flags().StringSliceP(cli.ArgNameLabel, "l", []string{}, `labels to attach to be included in the benchmark report`)
	rootCmd.Flags().BoolP(cli.ArgNameHeaders, "", true, `in tabular formats, whether to include headers in the report`)

	// Stdout
	rootCmd.Flags().BoolP(cli.ArgNamePipeStdout, "", true, `redirects external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP(cli.ArgNamePipeStderr, "", true, `redirects external commands standard error to benchy's standard error`)
	rootCmd.Flags().BoolP(cli.ArgNameDebug, "d", false, `logs extra debug information`)
	rootCmd.Flags().BoolP(cli.ArgNameSilent, "s", false, `logs only fatal errors`)

	_ = rootCmd.MarkFlagRequired(cli.ArgNameConfig)
	_ = rootCmd.MarkFlagFilename(cli.ArgNameConfig, "yml", "yaml", "json")
	_ = rootCmd.MarkFlagFilename(cli.ArgNameOutputFile, "txt", "csv", "md")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	// Subcommands
	rootCmd.AddCommand(subcmd.CreateConfigCommand())
	rootCmd.AddCommand(subcmd.CreateUpdateCommand(Version, ProgramName))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
