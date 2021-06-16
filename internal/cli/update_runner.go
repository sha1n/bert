package cli

import (
	"fmt"

	"github.com/sha1n/bert/api"
	clibcmd "github.com/sha1n/clib/pkg/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	gitHubRepoOwner = "sha1n"
	gitHubRepoName  = "bert"
)

// CreateUpdateCommand creates the 'config' sub command
func CreateUpdateCommand(version, binaryName string, ctx api.IOContext) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update",
		Long: fmt.Sprintf(`Checks if a newer release is available on GitHub and updates if so (see: https://github.com/%s/%s/releases). 
If '--tag' is specified, tries to update to the specified release tag regardless of whether it is more recent or not.`, gitHubRepoOwner, gitHubRepoName),
		Short: fmt.Sprintf(`Checks for and attempts to update %s to the latest or requested version`, binaryName),
		Run:   clibcmd.RunSelfUpdateFn(gitHubRepoOwner, gitHubRepoName, version, binaryName),
		PreRun: func(cmd *cobra.Command, args []string) {
			configureOutput(cmd, log.InfoLevel, ctx)
		},
	}

	cmd.Flags().String("tag", "", `the version tag to update to`)

	return cmd
}
