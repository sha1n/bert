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

type GetLatestReleaseFn = func() (github.Release, error)
type ResolveBinaryPathFn = func() (string, error)

func RunSelfUpdateFor(version, binaryName string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		CheckFatal(runSelfUpdateWith(version, binaryName, os.Executable, github.GetLatestRelease))
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
