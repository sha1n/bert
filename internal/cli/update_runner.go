package cli

import (
	"os"

	clibcmd "github.com/sha1n/clib/pkg/cmd"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	gitHusRepoOwner = "sha1n"
	gitHusRepoName  = "benchy"
)

// CreateUpdateCommand creates the 'config' sub command
func CreateUpdateCommand(version, binaryName string, ctx IOContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Long:  `Checks for a newer release on GitHub and updates if one is found (https://github.com/sha1n/benchy/releases)`,
		Short: `Checks for a newer release on GitHub and updates if one is found`,
		Run:   runSelfUpdateFn(version, binaryName, ctx),
	}

	return cmd
}

// runSelfUpdateFn runs the self update command based on the current version and binary name.
// currentVersion is used to determine whether a newer one is available
func runSelfUpdateFn(currentVersion, binaryName string, ctx IOContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// configureOutput(cmd, ctx)
		teardown := configureRichOutput(cmd, ctx)
		defer teardown()

		spinner := termite.NewDefaultSpinner()
		cancel, _ := spinner.Start()
		defer cancel()

		CheckFatal(clibcmd.RunSelfUpdate(gitHusRepoOwner, gitHusRepoName, currentVersion, binaryName, os.Executable, clibcmd.GetLatestRelease))

		log.Info("Done!")
	}
}
