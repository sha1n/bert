package subcmd

import (
	"github.com/sha1n/benchy/internal/cli"
	"github.com/spf13/cobra"
)

// CreateConfigCommand creates the 'config' sub command
func CreateUpdateCommand(version, binaryName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Long:  `Checks for a newer release on GitHub and updates if one is found (https://github.com/sha1n/benchy/releases)`,
		Short: `Checks for a newer release on GitHub and updates if one is found`,
		Run:   cli.RunSelfUpdateFor(version, binaryName),
	}

	return cmd
}
