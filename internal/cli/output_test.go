package cli

import (
	"os"
	"testing"

	"github.com/sha1n/benchy/test"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	cmd := aCommandWithArgs()
	configureOutput(cmd)

	assert.Equal(t, log.InfoLevel, log.StandardLogger().Level)
	assert.Equal(t, os.Stdout, log.StandardLogger().Out)
}

func TestDebugOn(t *testing.T) {
	cmd := aCommandWithArgs("-d")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd)

		assert.Equal(t, log.DebugLevel, log.StandardLogger().Level)
		assert.Equal(t, os.Stdout, log.StandardLogger().Out)
	}
	cmd.Execute()
}

func TestSilentOn(t *testing.T) {
	cmd := aCommandWithArgs("-s")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd)

		assert.Equal(t, log.PanicLevel, log.StandardLogger().Level)
		assert.Equal(t, os.Stderr, log.StandardLogger().Out)
	}
	cmd.Execute()
}

func aCommandWithArgs(args ...string) *cobra.Command {
	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString())
	rootCmd.SetArgs(append(args, "--config=../../test/data/integration.yaml"))

	return rootCmd
}
