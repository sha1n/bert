package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/sha1n/benchy/test"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

var expectedGoVersionOutput string
var itConfigFileArgValue = "--config=../../test/data/integration.yaml"

func init() {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	cmd := exec.Command("go", "version")
	cmd.Stdout = writer

	_ = cmd.Run()
	expectedGoVersionOutput = buf.String()
}

func TestBasic(t *testing.T) {
	output, err := runBenchmarkCommandWithPipedStdout(t, itConfigFileArgValue)

	assert.NoError(t, err)
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicMd(t *testing.T) {
	output, err := runBenchmarkCommandWithPipedStdout(t, itConfigFileArgValue, "--format=md")

	assert.NoError(t, err)
	assert.Contains(t, output, "|NAME|")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicMdRaw(t *testing.T) {
	output, err := runBenchmarkCommandWithPipedStdout(t, itConfigFileArgValue, "--format=md/raw")

	assert.NoError(t, err)
	assert.Contains(t, output, "|NAME|")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicCsv(t *testing.T) {
	output, err := runBenchmarkCommandWithPipedStdout(t, itConfigFileArgValue, "--format=csv")

	assert.NoError(t, err)
	assert.Contains(t, output, ",NAME,")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicCsvRaw(t *testing.T) {
	output, err := runBenchmarkCommandWithPipedStdout(t, itConfigFileArgValue, "--format=csv/raw")

	assert.NoError(t, err)
	assert.Contains(t, output, ",NAME,")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestWithMissingConfigFile(t *testing.T) {
	nonExistingConfigArg := fmt.Sprintf("-c=/tmp/%s", test.RandomString())
	_, _ = runBenchmarkCommandExpectPanic(t, nonExistingConfigArg)
}

func TestWithInvalidConfigFile(t *testing.T) {
	invalidConfig := "-c=../../test/data/invalid_config.yml"
	_, _ = runBenchmarkCommandExpectPanic(t, invalidConfig)
}

func TestWithCombinedDebugAndSilent(t *testing.T) {
	_, _ = runBenchmarkCommandExpectPanic(t, "-s", "-d", itConfigFileArgValue)
}

func runBenchmarkCommandWithPipedStdout(t *testing.T, args ...string) (output string, err error) {
	defer expectNoPanic(t)

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	originalWriter := log.StandardLogger().Out
	log.StandardLogger().SetOutput(buf)
	defer log.StandardLogger().SetOutput(originalWriter)

	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString())
	rootCmd.SetArgs(append(args, "--pipe-stdout"))
	rootCmd.SetOut(writer)
	rootCmd.SetErr(os.Stderr)

	err = rootCmd.Execute()

	return buf.String(), err
}

func runBenchmarkCommandExpectPanic(t *testing.T, args ...string) (output string, err error) {
	defer expectPanicWithError(t)

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	originalWriter := log.StandardLogger().Out
	log.StandardLogger().SetOutput(buf)
	defer log.StandardLogger().SetOutput(originalWriter)

	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString())
	rootCmd.SetArgs(args)
	rootCmd.SetOut(writer)
	rootCmd.SetErr(os.Stderr)

	err = rootCmd.Execute()

	return buf.String(), err
}

func expectNoPanic(t *testing.T) {
	if o := recover(); o != nil {
		if err, ok := o.(error); ok {
			assert.NoError(t, err)
		} else {
			assert.NoError(t, fmt.Errorf("%v", o))
		}
	}
}

func expectPanicWithError(t *testing.T) {
	if o := recover(); o != nil {
		if err, ok := o.(error); ok {
			assert.Error(t, err)
		} else {
			assert.Fail(t, fmt.Sprintf("A panic with an error was expected, but got: %v", o))
		}
	} else {
		assert.Fail(t, "A panic with an error was expected, but got nothing")
	}

}
