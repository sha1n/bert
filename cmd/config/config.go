package config

import (
	"github.com/sha1n/benchy/internal/cli"
	"github.com/spf13/cobra"
)

// CreateConfigCommand creates the 'config' sub command
func CreateConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Long:  `Interactively walks through a benchmark configuration creation process`,
		Short: `interactively creates a benchmark config`,
		Run:   cli.CreateConfig,
	}

	cmd.Flags().StringP(cli.ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)

	_ = cmd.MarkFlagFilename(cli.ArgNameOutputFile, "yml", "yaml")

	return cmd
}
