package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg"
	"github.com/sha1n/gommons/pkg/test"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var expectedGoVersionOutput string
var itConfigFilePath = "../../test/data/integration.yaml"
var itConfigFileArgValue = fmt.Sprintf("--config=%s", itConfigFilePath)

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
		itConfigFileArgValue, "--format=csv",
	)
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
	nonExistingConfigArg := fmt.Sprintf("-c=/tmp/%s", gommonstest.RandomString())
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, nonExistingConfigArg)
}

func TestWithInvalidConfigFile(t *testing.T) {
	invalidConfig := "-c=../../test/data/invalid_config.yml"
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, invalidConfig)
}

func TestWithCombinedDebugAndSilent(t *testing.T) {
	runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t, "-s", "-d", itConfigFileArgValue)
}

func Test_loadSpecWithExecutionsOverride(t *testing.T) {
	expectedSpec, _ := pkg.LoadSpec(itConfigFilePath)
	expectedSpec.Executions = expectedSpec.Executions + rand.Intn(10)
	command := newDummyCommandWith("-c", itConfigFilePath, "--executions", fmt.Sprint(expectedSpec.Executions))

	spec, err := loadSpec(command, []string{})

	assert.NoError(t, err)
	assert.Equal(t, expectedSpec, spec)
}

func Test_loadSpecFromPositionalArguments(t *testing.T) {
	expectedExecutions := rand.Intn(100)
	expectedValidSpec, _ := pkg.CreateSpecFrom(
		expectedExecutions,
		false,
		api.CommandSpec{Cmd: []string{"cmd", "-a"}},
		api.CommandSpec{Cmd: []string{"cmd", "-b"}},
		api.CommandSpec{Cmd: []string{"cmd", "a string"}},
	)

	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name     string
		args     args
		wantSpec api.BenchmarkSpec
		wantErr  bool
	}{
		{
			name: "call with positional single-quote framed commands",
			args: args{
				cmd:  newDummyCommandWith("--executions", fmt.Sprint(expectedExecutions)),
				args: []string{"'cmd -a'", "'cmd -b'", "'cmd \"a string\"'"},
			},
			wantSpec: expectedValidSpec,
			wantErr:  false,
		},
		{
			name: "call with positional double-quote framed commands",
			args: args{
				cmd:  newDummyCommandWith("--executions", fmt.Sprint(expectedExecutions)),
				args: []string{"\"cmd -a\"", "\"cmd -b\"", "\"cmd 'a string'\""},
			},
			wantSpec: expectedValidSpec,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpec, err := loadSpec(tt.args.cmd, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, expectedValidSpec, gotSpec)
		})
	}
}

func Test_validatePositionalArgs(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "call with no positional and config path",
			args: args{
				cmd:  newDummyCommandWith("-c", "/some-file"),
				args: []string{},
			},
			wantErr: false,
		},
		{
			name: "call with no positional and no config path",
			args: args{
				cmd:  newDummyCommandWith(),
				args: []string{},
			},
			wantErr: true,
		},
		{
			name: "call with positional and no executions param",
			args: args{
				cmd:  newDummyCommandWith("-c", "/some-file"),
				args: []string{test.RandomString(), test.RandomString()},
			},
			wantErr: true,
		},
		{
			name: "call with positional and invalid executions param",
			args: args{
				cmd:  newDummyCommandWith("--executions", "-1"),
				args: []string{test.RandomString(), test.RandomString()},
			},
			wantErr: true,
		},
		{
			name: "call with positional and valid executions param",
			args: args{
				cmd:  newDummyCommandWith("--executions", "100"),
				args: []string{test.RandomString(), test.RandomString()},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePositionalArgs(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validatePositionalArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func runBenchmarkCommandWithPipedStdoutputsAnd(t *testing.T, assert func(stdout, stderr string, err error), args ...string) {
	defer expectNoPanic(t)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	ioContext := api.NewIOContext()
	ioContext.StdoutWriter = outBuf
	ioContext.StderrWriter = errBuf
	rootCmd := NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--pipe-stdout=true", "--pipe-stderr=true"))
	rootCmd.SetOut(ioContext.StdoutWriter)
	rootCmd.SetErr(ioContext.StderrWriter)

	err := rootCmd.Execute()

	assert(outBuf.String(), errBuf.String(), err)
}

func runBenchmarkCommandWithPipedStdoutAndExpectPanicWith(t *testing.T, args ...string) {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	ioContext := api.NewIOContext()
	ioContext.StdoutWriter = bufio.NewWriter(buf)
	ioContext.StderrWriter = bufio.NewWriter(buf)

	rootCmd := NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), ioContext)
	rootCmd.SetArgs(append(args, "--pipe-stdout=true"))
	rootCmd.SetOut(writer)
	rootCmd.SetErr(os.Stderr)

	assert.Panics(t, func() { _ = rootCmd.Execute() })
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

func newDummyCommandWith(args ...string) *cobra.Command {
	rootCmd := NewRootCommand(gommonstest.RandomString(), gommonstest.RandomString(), gommonstest.RandomString(), api.NewIOContext())
	rootCmd.SetArgs(args)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {}

	c, _ := rootCmd.ExecuteC()

	return c
}
