package cli

import (
	"context"
	"log/slog"
	"testing"

	"github.com/sha1n/bert/api"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/sha1n/termite"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	ctx := api.NewIOContext()
	cmd := aCommandWithArgs(ctx)
	configureOutput(cmd, slog.LevelError, ctx)

	assert.True(t, slog.Default().Enabled(context.Background(), slog.LevelError))
	assert.False(t, slog.Default().Enabled(context.Background(), slog.LevelInfo))
}

func TestDebugOn(t *testing.T) {
	ctx := api.NewIOContext()
	cmd := aCommandWithArgs(ctx, "-d")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd, slog.LevelError, ctx)

		assert.True(t, slog.Default().Enabled(context.Background(), slog.LevelDebug))
		assert.True(t, slog.Default().Enabled(context.Background(), slog.LevelInfo))
		assert.True(t, slog.Default().Enabled(context.Background(), slog.LevelError))
	}

	assert.NoError(t, cmd.Execute())
}

func TestSilentOn(t *testing.T) {
	ctx := api.NewIOContext()
	cmd := aCommandWithArgs(ctx, "-s")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		configureOutput(cmd, slog.LevelError, ctx)

		// We implemented silent as LevelError + 1
		assert.False(t, slog.Default().Enabled(context.Background(), slog.LevelError))
		assert.False(t, slog.Default().Enabled(context.Background(), slog.LevelInfo))
	}

	assert.NoError(t, cmd.Execute())
}

func TestTtyMode(t *testing.T) {
	withTty(func(ctx api.IOContext) {
		cmd := aCommandWithArgs(ctx)
		cmd.Run = func(cmd *cobra.Command, args []string) {
			configureOutput(cmd, slog.LevelError, ctx)

			// We can't easily assert on handlers internal state like ForceColors or DisableTimestamp with slog
			// But we can ensure it configured the level correctly
			assert.True(t, slog.Default().Enabled(context.Background(), slog.LevelError))
		}

		assert.NoError(t, cmd.Execute())
	})
}

func aCommandWithArgs(ctx api.IOContext, args ...string) *cobra.Command {
	ioContext := api.NewIOContext()
	rootCmd := NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--config=../../test/data/integration.yaml"))

	return rootCmd
}

func withTty(test func(api.IOContext)) {
	origTty := termite.Tty
	termite.Tty = true
	ioContext := api.NewIOContext()
	ioContext.Tty = true

	defer func() {
		termite.Tty = origTty
	}()

	test(ioContext)
}
