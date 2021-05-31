package cli

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sha1n/benchy/internal/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// GetLatestReleaseFn ...
type GetLatestReleaseFn = func() (github.Release, error)

// ResolveBinaryPathFn ...
type ResolveBinaryPathFn = func() (string, error)

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
		defer configureNonInteractiveOutput(cmd, ctx)()

		CheckFatal(runSelfUpdateWith(currentVersion, binaryName, os.Executable, github.GetLatestRelease))

		log.Info("Done!")
	}
}

func runSelfUpdateWith(version, binaryName string, resolveBinaryPathFn ResolveBinaryPathFn, getLatestReleaseFn GetLatestReleaseFn) (err error) {
	var binaryPath string
	if binaryPath, err = resolveBinaryPathFn(); err != nil {
		return err
	}

	log.Infof("Fetching latest release...")
	var release github.Release
	if release, err = getLatestReleaseFn(); err != nil {
		return err
	}

	tagName := release.TagName()
	log.Infof("Latest release tag is %s", tagName)
	log.Infof("Current version is %s", version)

	if tagName != "" && tagName != version && semver.Compare(tagName, version) > 0 {
		log.Infof("Downloading version %s...", tagName)
		var rc io.ReadCloser
		if rc, err = release.DownloadAsset(); err != nil {
			return err
		}

		var content []byte
		if content, err = ioutil.ReadAll(rc); err == nil {
			return ioutil.WriteFile(binaryPath, content, 0755)
		}

	} else {
		log.Infof("You are already running the latest version of %s!", binaryName)
	}

	return err
}
