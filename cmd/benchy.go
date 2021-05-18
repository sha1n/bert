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

	rootCmd.Flags().StringP(cli.ArgNameConfig, "c", "", `config file path`)
	rootCmd.Flags().StringP(cli.ArgNameOutputFile, "o", "", `output file path`)
	rootCmd.Flags().StringP(cli.ArgNameFormat, "f", "txt", `summary format. One of: 'txt', 'csv' (default: txt)`)
	rootCmd.Flags().BoolP(cli.ArgNamePipeStdout, "", true, `redirects external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP(cli.ArgNamePipeStderr, "", true, `redirects external commands standard error to benchy's standard error`)
	rootCmd.Flags().BoolP(cli.ArgNameDebug, "d", false, `logs extra debug information`)
	rootCmd.Flags().StringSliceP(cli.ArgNameLabel, "l", []string{}, `labels to attach to this benchmark.`)

	cobra.MarkFlagRequired(rootCmd.Flags(), "config")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	_ = rootCmd.Execute()
}
