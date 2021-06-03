package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	clibtest "github.com/sha1n/clib/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var (
	userInputExecutions          = clibtest.RandomUint()
	userInputAlternate           = clibtest.RandomBool()
	userInputScenarioName        = clibtest.RandomString()
	userInputScenarioWorkingDir  = os.TempDir() // has to exist
	userInputDefineEnvVarsAnswer = "y"
	userInputEnvVarValue         = fmt.Sprintf("X=%s", clibtest.RandomString())
	userInputCommand             = fmt.Sprintf("cmd %s", clibtest.RandomString())
)

func TestBasicInteractiveFlow(t *testing.T) {
	ctx := givenIOContextWithInputContent(userInput())

	rootCmd, configPath, teardown := configureCommand(t, ctx)
	defer teardown()

	err := rootCmd.Execute()
	assert.NoError(t, err)

	actual, err := pkg.LoadSpec(configPath)

	assert.NoError(t, err)
	assert.Equal(t, expectedSpec(), actual)
}

func configureCommand(t *testing.T, ctx IOContext) (command *cobra.Command, configPath string, teardown func()) {
	rootCmd := NewRootCommand(clibtest.RandomString(), clibtest.RandomString(), clibtest.RandomString(), ctx)
	cmd := CreateConfigCommand(ctx)
	rootCmd.AddCommand(cmd)

	tmpFile, err := ioutil.TempFile("", "TestBasicInteractiveFlow")

	assert.NoError(t, err)

	args := []string{"config", "--out-file", tmpFile.Name()}
	// cmd.SetArgs(args)

	assert.NoError(t, cmd.Flags().Set("out-file", tmpFile.Name()))
	rootCmd.SetArgs(args)

	return rootCmd, tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
}

func userInput() string {
	return fmt.Sprintf(`%d
%v
%s
%s
%s
%s












%s


`,
		userInputExecutions, userInputAlternate, userInputScenarioName, userInputScenarioWorkingDir, userInputDefineEnvVarsAnswer, userInputEnvVarValue, userInputCommand)
}

func expectedSpec() api.BenchmarkSpec {
	kv := strings.Split(userInputEnvVarValue, "=")

	return api.BenchmarkSpec{
		Executions: int(userInputExecutions),
		Alternate:  userInputAlternate,
		Scenarios: []api.ScenarioSpec{
			{
				Name:             userInputScenarioName,
				WorkingDirectory: userInputScenarioWorkingDir,
				Env:              map[string]string{kv[0]: kv[1]},
				BeforeAll:        nil,
				AfterAll:         nil,
				BeforeEach:       nil,
				AfterEach:        nil,
				Command: &api.CommandSpec{
					Cmd: strings.Split(userInputCommand, " "),
				},
			},
		},
	}
}
