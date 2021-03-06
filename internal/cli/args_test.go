package cli

import (
	"fmt"
	"testing"

	"github.com/sha1n/bert/api"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestIsExperimentEnabledWithNoExperimentalFlag(t *testing.T) {
	withCommandWithArgs(t, func(cmd *cobra.Command) {
		assert.False(t, IsExperimentEnabled(cmd, gommonstest.RandomString()))
	})
}

func TestIsExperimentEnabledWithExperimentalFlagAndNoMatch(t *testing.T) {
	flagValue := experimentalFlagWith(gommonstest.RandomString())

	withCommandWithArgs(
		t,
		func(cmd *cobra.Command) {
			assert.False(t, IsExperimentEnabled(cmd, gommonstest.RandomString()))
		},
		flagValue,
	)
}

func TestIsExperimentEnabledWithMatchingExperimentalFlag(t *testing.T) {
	featureName := gommonstest.RandomString()
	flagValue := experimentalFlagWith(featureName)

	withCommandWithArgs(
		t,
		func(cmd *cobra.Command) {
			assert.True(t, IsExperimentEnabled(cmd, featureName))
		},
		flagValue,
	)
}

func experimentalFlagWith(value string) string {
	return fmt.Sprintf("--%s=%s", ArgNameExperimental, value)
}

func withCommandWithArgs(t *testing.T, doTest func(cmd *cobra.Command), args ...string) {
	ctx := api.NewIOContext()
	cmd := NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ctx)
	cmd.SetArgs(append(args, "--config=xxx"))
	cmd.Run = func(c *cobra.Command, args []string) {
		doTest(c)
	}

	assert.NoError(t, cmd.Execute())
}
