package cli

import (
	"testing"

	"github.com/sha1n/benchy/test"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	cmd := aCommandWithArgs()
	configureOutput(cmd)

	assert.Equal(t, log.InfoLevel, log.StandardLogger().Level)
	assert.Equal(t, StdoutWriter, log.StandardLogger().Out)
}

func TestDebugOn(t *testing.T) {
	cmd := aCommandWithArgs("-d")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd)

		assert.Equal(t, log.DebugLevel, log.StandardLogger().Level)
		assert.Equal(t, StdoutWriter, log.StandardLogger().Out)
	}
	cmd.Execute()
}

func TestSilentOn(t *testing.T) {
	cmd := aCommandWithArgs("-s")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd)

		assert.Equal(t, log.PanicLevel, log.StandardLogger().Level)
		assert.Equal(t, StderrWriter, log.StandardLogger().Out)
	}
	cmd.Execute()
}

func TestTTYModeConfiguration(t *testing.T) {
	cmd := aCommandWithArgs("-s")
	origTty := termite.Tty
	termite.Tty = true
	defer func() {
		termite.Tty = origTty
	}()

	cancel := configureNonInteractiveOutput(cmd)
	assert.NotEqual(t, StdoutWriter, log.StandardLogger().Out)
	assert.IsType(t, &alwaysRewritingWriter{}, log.StandardLogger().Out)
	cancel()
}

func aCommandWithArgs(args ...string) *cobra.Command {
	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString())
	rootCmd.SetArgs(append(args, "--config=../../test/data/integration.yaml"))

	return rootCmd
}
