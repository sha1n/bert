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
	ctx := NewIOContext()
	cmd := aCommandWithArgs(ctx)
	configureOutput(cmd, ctx)

	assert.Equal(t, log.InfoLevel, log.StandardLogger().Level)
	assert.Equal(t, ctx.StderrWriter, log.StandardLogger().Out)
}

func TestDebugOn(t *testing.T) {
	ctx := NewIOContext()
	cmd := aCommandWithArgs(ctx, "-d")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd, ctx)

		assert.Equal(t, log.DebugLevel, log.StandardLogger().Level)
		assert.Equal(t, ctx.StderrWriter, log.StandardLogger().Out)
	}

	assert.NoError(t, cmd.Execute())
}

func TestSilentOn(t *testing.T) {
	ctx := NewIOContext()
	cmd := aCommandWithArgs(ctx, "-s")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd, ctx)

		assert.Equal(t, log.PanicLevel, log.StandardLogger().Level)
		assert.Equal(t, ctx.StderrWriter, log.StandardLogger().Out)
	}

	assert.NoError(t, cmd.Execute())
}

func TestTtyMode(t *testing.T) {
	withTty(func(ctx IOContext) {
		cmd := aCommandWithArgs(ctx)
		cmd.Run = func(cmd *cobra.Command, args []string) {
			configureOutput(cmd, ctx)

			assert.Equal(t, ctx.StderrWriter, log.StandardLogger().Out)
			assert.True(t, log.StandardLogger().Formatter.(*log.TextFormatter).DisableTimestamp)
			assert.True(t, log.StandardLogger().Formatter.(*log.TextFormatter).ForceColors)
		}

		assert.NoError(t, cmd.Execute())
	})
}

func aCommandWithArgs(ctx IOContext, args ...string) *cobra.Command {
	ioContext := NewIOContext()
	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--config=../../test/data/integration.yaml"))

	return rootCmd
}

func withTty(test func(IOContext)) {
	origTty := termite.Tty
	termite.Tty = true
	ioContext := NewIOContext()
	ioContext.Tty = true

	defer func() {
		termite.Tty = origTty
	}()

	test(ioContext)
}
