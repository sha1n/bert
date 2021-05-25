package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	"github.com/sha1n/benchy/test"
	"github.com/stretchr/testify/assert"
)

var (
	userInputExecutions          = test.RandomUint()
	userInputAlternate           = test.RandomBool()
	userInputScenarioName        = test.RandomString()
	userInputScenarioWorkingDir  = os.TempDir() // has to exist
	userInputDefineEnvVarsAnswer = "y"
	userInputEnvVarValue         = fmt.Sprintf("X=%s", test.RandomString())
	userInputCommand             = fmt.Sprintf("cmd %s", test.RandomString())
)

func TestBasicInteractiveFlow(t *testing.T) {
	teardown := givenStdInWith(userInput())
	defer teardown()

	cmd := CreateConfigCommand()

	tmpFile, err := ioutil.TempFile("", "TestBasicInteractiveFlow")
	defer os.Remove(tmpFile.Name())

	assert.NoError(t, err)

	args := []string{"--out-file", tmpFile.Name()}
	cmd.SetArgs(args)
	assert.NoError(t, cmd.Flags().Set("out-file", tmpFile.Name()))

	CreateConfig(cmd, args)

	actual, err := pkg.LoadSpec(tmpFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, expectedSpec(), actual)
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

func expectedSpec() *api.BenchmarkSpec {
	kv := strings.Split(userInputEnvVarValue, "=")

	return &api.BenchmarkSpec{
		Executions: int(userInputExecutions),
		Alternate:  userInputAlternate,
		Scenarios: []*api.ScenarioSpec{
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
