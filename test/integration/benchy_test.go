package integration

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedGoVersionOutput string

func init() {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	cmd := exec.Command("go", "version")
	cmd.Stdout = writer

	_ = cmd.Run()
	expectedGoVersionOutput = buf.String()
}

func TestBasic(t *testing.T) {
	output, err := runBenchmarkCommand("go", "run", "../../cmd/main.go", "-s", "--config=../data/integration.yaml")

	assert.NoError(t, err)
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestBasicMd(t *testing.T) {
	output, err := runBenchmarkCommand("go", "run", "../../cmd/main.go", "-s", "--config=../data/integration.yaml", "--format=md")

	assert.NoError(t, err)
	assert.Contains(t, output, "|NAME|")
	assert.Contains(t, output, expectedGoVersionOutput)
}
func TestBasicCsv(t *testing.T) {
	output, err := runBenchmarkCommand("go", "run", "../../cmd/main.go", "-s", "--config=../data/integration.yaml", "--format=csv")

	assert.NoError(t, err)
	assert.Contains(t, output, ",NAME,")
	assert.Contains(t, output, expectedGoVersionOutput)
}

func TestExit1WhenRequiredArgsAreMissing(t *testing.T) {
	cmd := exec.Command("go", "run", "cmd/main.go")

	assert.Error(t, cmd.Run())
}

func runBenchmarkCommand(cmd ...string) (string, error) {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdout = writer
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()

	return buf.String(), err
}
