package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"os/exec"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

func TestExecuteReturnsErrorOnCommandFailure(t *testing.T) {
	defaultWorkingDir := "default/dir"
	env := map[string]string{}
	exec := NewCommandExecutor(true, true)
	cmdSpec := aCommandSpec(aNonExistingCommand(), "")

	err := exec.Execute(cmdSpec, defaultWorkingDir, env)

	assert.Error(t, err)
}

func TestConfigureCommandWithEmptyWorkingDir(t *testing.T) {
	defaultWorkingDir := ""
	env := map[string]string{}
	cmd := aNonExistingCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, defaultWorkingDir, execCmd.Dir)
	assert.Equal(t, os.Environ(), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithTildeDefaultWorkingDir(t *testing.T) {
	defaultWorkingDir := "~"
	env := map[string]string{}
	cmd := aNonExistingCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, userHomeDir(), execCmd.Dir)
	assert.Equal(t, os.Environ(), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithTildeLocalWorkingDir(t *testing.T) {
	defaultWorkingDir := ""
	env := map[string]string{}
	cmd := aNonExistingCommand()

	execCmd := configureCommand(aCommandSpec(cmd, "~"), defaultWorkingDir, env)

	assert.Equal(t, userHomeDir(), execCmd.Dir)
	assert.Equal(t, os.Environ(), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithCustomEnv(t *testing.T) {
	defaultWorkingDir := ""
	env := aCustomEnv()
	cmd := aNonExistingCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, defaultWorkingDir, execCmd.Dir)
	assert.Equal(t, expectedEnvFor(env), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithStdoutPiping(t *testing.T) {
	execCmd := configureCommandWithIOSpec(true, false)

	assert.Equal(t, log.StandardLogger().Out, execCmd.Stdout)
	assert.Equal(t, nil, execCmd.Stderr)
}

func TestConfigureCommandWithStderrPiping(t *testing.T) {
	execCmd := configureCommandWithIOSpec(false, true)

	assert.Equal(t, nil, execCmd.Stdout)
	assert.Equal(t, log.StandardLogger().Out, execCmd.Stderr)
}

func configureCommandWithIOSpec(pipeStdout, pipeStderr bool) *exec.Cmd {
	spec := aCommandSpec(aNonExistingCommand(), "")
	executor := NewCommandExecutor(pipeStdout, pipeStderr).(*commandExecutor)
	execCmd := exec.Command(spec.Cmd[0], spec.Cmd[1:]...)

	executor.configureCommand(spec, execCmd, "", map[string]string{})

	return execCmd
}

func configureCommand(spec *api.CommandSpec, defaultWorkingDir string, env map[string]string) *exec.Cmd {
	executor := NewCommandExecutor(false, false).(*commandExecutor)
	execCmd := exec.Command(spec.Cmd[0], spec.Cmd[1:]...)

	executor.configureCommand(spec, execCmd, defaultWorkingDir, env)

	return execCmd
}

func aNonExistingCommand() []string {
	return []string{"dummyCmd", "-arg"}
}

func aCustomEnv() map[string]string {
	return map[string]string{"KEY": fmt.Sprintf("%d", time.Now().Nanosecond())}
}

func expectedEnvFor(e map[string]string) []string {
	return append(toEnvVarsArray(e), os.Environ()...) // user vars are expected to be first
}

func userHomeDir() string {
	if p, e := os.UserHomeDir(); e == nil {
		return p
	}

	return ""
}
func aCommandSpec(cmd []string, workingDir string) *api.CommandSpec {
	return &api.CommandSpec{
		WorkingDirectory: workingDir,
		Cmd:              cmd,
	}
}
