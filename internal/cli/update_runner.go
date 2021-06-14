package cli

import (
	"os"

	"github.com/sha1n/bert/api"
	clibcmd "github.com/sha1n/clib/pkg/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	gitHusRepoOwner = "sha1n"
	gitHusRepoName  = "bert"
)

// CreateUpdateCommand creates the 'config' sub command
func CreateUpdateCommand(version, binaryName string, ctx api.IOContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Long:  `Checks for a newer release on GitHub and updates if one is found (https://github.com/sha1n/bert/releases)`,
		Short: `Checks for a newer release on GitHub and updates if one is found`,
		Run:   runSelfUpdateFn(version, binaryName, ctx),
	}

	return cmd
}

// runSelfUpdateFn runs the self update command based on the current version and binary name.
// currentVersion is used to determine whether a newer one is available
func runSelfUpdateFn(currentVersion, binaryName string, ctx api.IOContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		configureOutput(cmd, log.InfoLevel, ctx)

		CheckFatal(clibcmd.RunSelfUpdate(gitHusRepoOwner, gitHusRepoName, currentVersion, binaryName, os.Executable, clibcmd.GetLatestRelease))

		log.Info("Done!")
	}
}
