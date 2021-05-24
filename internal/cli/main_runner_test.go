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
var mainSourceFilePath = "../../cmd/main.go"

func init() {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	cmd := exec.Command("go", "version")
	cmd.Stdout = writer

	_ = cmd.Run()
	expectedGoVersionOutput = buf.String()
}

func TestBasic(t *testing.T) {
	output, err := runBenchmarkCommandWith(t, "-s", itConfigFileArgValue)

	assert.NoError(t, err)
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicMd(t *testing.T) {
	output, err := runBenchmarkCommandWith(t, "-s", itConfigFileArgValue, "--format=md")

	assert.NoError(t, err)
	assert.Contains(t, output, "|NAME|")
	assert.Contains(t, output, expectedGoVersionOutput)
}
func TestBasicCsv(t *testing.T) {
	output, err := runBenchmarkCommandWith(t, "-s", itConfigFileArgValue, "--format=csv")

	assert.NoError(t, err)
	assert.Contains(t, output, ",NAME,")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func runBenchmarkCommandWith(t *testing.T, args ...string) (output string, err error) {
	defer func() {
		if o := recover(); o != nil {
			if err, ok := o.(error); ok {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, fmt.Errorf("%v", o))
			}
		}
	}()

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
