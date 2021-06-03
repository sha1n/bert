package cli

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	clibtest "github.com/sha1n/clib/pkg/test"

	"github.com/stretchr/testify/assert"
)

const (
	itConfigFileArgValue = "--config=../../test/data/integration.yaml"
	expectedScenarioName = "TEST_SCENARIO_NAME"
)

var (
	expectedScenarioMdCell  = fmt.Sprintf("|%s|", expectedScenarioName)
	expectedScenarioCsvCell = fmt.Sprintf(",%s,", expectedScenarioName)
	expectedGoVersionOutput string
)

func init() {
	buf := new(bytes.Buffer)

	cmd := exec.Command("go", "version")
	cmd.Stdout = buf

	_ = cmd.Run()
	expectedGoVersionOutput = buf.String()
}

func TestBasic(t *testing.T) {
	expectedStartupLogMessage := "Starting benchy..."
	runBenchmarkCommand(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)

			assert.Contains(t, stdout, expectedScenarioName)

			assert.Contains(t, stderr, expectedStartupLogMessage)
			assert.NotContains(t, stdout, expectedStartupLogMessage)

		},
		itConfigFileArgValue,
	)
}

func TestBasicWithPipedStdout(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputs(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, expectedScenarioName)
			assert.Contains(t, stderr, expectedGoVersionOutput)
			assert.NotContains(t, stdout, expectedGoVersionOutput)
		},
		itConfigFileArgValue,
	)

}

func TestBasicMd(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputs(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, expectedScenarioMdCell)
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=md",
	)
}

func TestBasicMdRaw(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputs(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, expectedScenarioMdCell)
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=md/raw",
	)

}

func TestBasicCsv(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputs(t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, expectedScenarioCsvCell)
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=csv")

}

func TestBasicCsvRaw(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutputs(
		t,
		func(stdout, stderr string, err error) {
			assert.NoError(t, err)
			assert.Contains(t, stdout, expectedScenarioCsvCell)
			assert.Contains(t, stderr, expectedGoVersionOutput)
		},
		itConfigFileArgValue, "--format=csv/raw",
	)

}

func TestWithMissingConfigFile(t *testing.T) {
	nonExistingConfigArg := fmt.Sprintf("-c=/tmp/%s", clibtest.RandomString())
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, nonExistingConfigArg)
}

func TestWithInvalidConfigFile(t *testing.T) {
	invalidConfig := "-c=../../test/data/invalid_config.yml"
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, invalidConfig)
}

func TestWithCombinedDebugAndSilent(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, "-s", "-d", itConfigFileArgValue)
}

func runBenchmarkCommand(t *testing.T, assert func(stdout, stderr string, err error), args ...string) {
	defer expectNoPanic(t)

	ctx, err := runBenchmark(args...)

	assert(ctx.StdoutWriter.(*bytes.Buffer).String(), ctx.StderrWriter.(*bytes.Buffer).String(), err)
}

func runBenchmarkCommandWithPipedStdoutputs(t *testing.T, assert func(stdout, stderr string, err error), args ...string) {
	runBenchmarkCommand(t, assert, append(args, "--pipe-stdout", "--pipe-stderr")...)
}

func runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t *testing.T, args ...string) {
	defer expectPanicWithError(t)

	_, _ = runBenchmark(args...)
}

func runBenchmark(args ...string) (IOContext, error) {
	ctx := newIOContext()
	rootCmd := NewRootCommand(clibtest.RandomString(), clibtest.RandomString(), clibtest.RandomString(), ctx)
	rootCmd.SetArgs(args)
	rootCmd.SetOut(ctx.StdoutWriter)
	rootCmd.SetErr(ctx.StderrWriter)

	return ctx, rootCmd.Execute()

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

func newIOContext() IOContext {
	ctx := NewIOContext()
	ctx.Tty = clibtest.RandomBool()
	ctx.StdoutWriter = new(bytes.Buffer)
	ctx.StderrWriter = new(bytes.Buffer)

	return ctx
}
