package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/sha1n/benchy/test"

	"github.com/stretchr/testify/assert"
)

var expectedGoVersionOutput string
var itConfigFileArgValue = "--config=../../test/data/integration.yaml"

func init() {
	buf := new(bytes.Buffer)

	cmd := exec.Command("go", "version")
	cmd.Stdout = buf

	_ = cmd.Run()
	expectedGoVersionOutput = buf.String()
}

func TestBasic(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputsAnd(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, "NAME")
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue,
	)

}

func TestBasicMd(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputsAnd(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, "|NAME|")
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=md",
	)
}

func TestBasicMdRaw(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputsAnd(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, "|NAME|")
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=md/raw",
	)

}

func TestBasicCsv(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputsAnd(t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, ",NAME,")
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=csv")

}

func TestBasicCsvRaw(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputsAnd(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, ",NAME,")
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=csv/raw",
	)

}

func TestWithMissingConfigFile(t *testing.T) {
	nonExistingConfigArg := fmt.Sprintf("-c=/tmp/%s", test.RandomString())
	_, _ = runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, nonExistingConfigArg)
}

func TestWithInvalidConfigFile(t *testing.T) {
	invalidConfig := "-c=../../test/data/invalid_config.yml"
	_, _ = runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, invalidConfig)
}

func TestWithCombinedDebugAndSilent(t *testing.T) {
	_, _ = runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, "-s", "-d", itConfigFileArgValue)
}

func runBenchmarkCommandWithPipedStdoutputsAnd(t *testing.T, assert func(stdout, stderr string, err error), args ...string) {
	defer expectNoPanic(t)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	ioContext := NewIOContext()
	ioContext.StdoutWriter = outBuf
	ioContext.StderrWriter = errBuf
	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--pipe-stdout=true", "--pipe-stderr=true"))
	rootCmd.SetOut(ioContext.StdoutWriter)
	rootCmd.SetErr(ioContext.StderrWriter)

	err := rootCmd.Execute()

	assert(outBuf.String(), errBuf.String(), err)
}

func runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t *testing.T, args ...string) (output string, err error) {
	defer expectPanicWithError(t)

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	ioContext := NewIOContext()
	ioContext.StdoutWriter = bufio.NewWriter(buf)
	ioContext.StderrWriter = bufio.NewWriter(buf)

	rootCmd := NewRootCommand(test.RandomString(), test.RandomString(), test.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--pipe-stdout=true"))
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
