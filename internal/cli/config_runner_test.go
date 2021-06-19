package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var (
	userInputExecutions          = gommonstest.RandomUint()
	userInputAlternate           = gommonstest.RandomBool()
	userInputScenarioName        = gommonstest.RandomString()
	userInputScenarioWorkingDir  = os.TempDir() // has to exist
	userInputDefineEnvVarsAnswer = "y"
	userInputEnvVarValue         = fmt.Sprintf("X=%s", gommonstest.RandomString())
	userInputCommand             = fmt.Sprintf("cmd %s", gommonstest.RandomString())
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

// Making sure the example we provide to the user is valid
func TestExampleSpecValidity(t *testing.T) {
	ctx := api.NewIOContext()
	buffer := new(bytes.Buffer)
	ctx.StdoutWriter = buffer

	rootCmd := configureExampleCommand(t, ctx)

	err := rootCmd.Execute()
	assert.NoError(t, err)

	actual, err := pkg.LoadSpecFromYamlData(buffer.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestExampleOutFile(t *testing.T) {
	ctx := api.NewIOContext()
	rootCmd, configPath, teardown := configureExampleCommandWithOutFile(t, ctx)
	teardown()

	expected, err := pkg.LoadSpecFromYamlData([]byte(getExampleSpec()))
	assert.NoError(t, err)
	assert.NotNil(t, expected)

	err = rootCmd.Execute()
	assert.NoError(t, err)

	actual, err := pkg.LoadSpec(configPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func configureExampleCommand(t *testing.T, ctx api.IOContext) (rootCmd *cobra.Command) {
	rootCmd = NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ctx)
	configCmd := CreateConfigCommand(ctx)
	rootCmd.AddCommand(configCmd)

	rootCmd.SetArgs([]string{"config", "--example"})

	return
}

func configureExampleCommandWithOutFile(t *testing.T, ctx api.IOContext) (rootCmd *cobra.Command, configPath string, teardown func()) {
	rootCmd = NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ctx)
	configCmd := CreateConfigCommand(ctx)
	rootCmd.AddCommand(configCmd)

	tmpFile, err := ioutil.TempFile("", "configureCommand")
	assert.NoError(t, err)

	configPath = tmpFile.Name()

	args := []string{"config", "--example", "-o", configPath}

	rootCmd.SetArgs(args)

	return rootCmd, configPath, func() { os.Remove(tmpFile.Name()) }
}

func configureCommand(t *testing.T, ctx api.IOContext) (rootCmd *cobra.Command, configPath string, teardown func()) {
	rootCmd = NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ctx)
	cmd := CreateConfigCommand(ctx)
	rootCmd.AddCommand(cmd)

	tmpFile, err := ioutil.TempFile("", "configureCommand")
	assert.NoError(t, err)

	rootCmd.SetArgs([]string{"config", "--out-file", tmpFile.Name()})

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
