package pkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"os/exec"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCommandWithEmptyWorkingDir(t *testing.T) {
	defaultWorkingDir := ""
	env := map[string]string{}
	cmd := aCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, defaultWorkingDir, execCmd.Dir)
	assert.Equal(t, []string(nil), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithTildeDefaultWorkingDir(t *testing.T) {
	defaultWorkingDir := "~"
	env := map[string]string{}
	cmd := aCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, userHomeDir(), execCmd.Dir)
	assert.Equal(t, []string(nil), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithTildeLocalWorkingDir(t *testing.T) {
	defaultWorkingDir := ""
	env := map[string]string{}
	cmd := aCommand()

	execCmd := configureCommand(aCommandSpec(cmd, "~"), defaultWorkingDir, env)

	assert.Equal(t, userHomeDir(), execCmd.Dir)
	assert.Equal(t, []string(nil), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func TestConfigureCommandWithCustomEnv(t *testing.T) {
	defaultWorkingDir := ""
	env := aCustomEnv()
	cmd := aCommand()

	execCmd := configureCommand(aCommandSpec(cmd, ""), defaultWorkingDir, env)

	assert.Equal(t, defaultWorkingDir, execCmd.Dir)
	assert.Equal(t, expectedEnvFor(env), execCmd.Env)
	assert.Equal(t, cmd, execCmd.Args)
}

func configureCommand(spec *api.CommandSpec, defaultWorkingDir string, env map[string]string) *exec.Cmd {
	executor := NewCommandExecutor(false, false).(*commandExecutor)
	execCmd := exec.Command(spec.Cmd[0], spec.Cmd[1:]...)

	executor.configureCommand(spec, execCmd, defaultWorkingDir, env)

	return execCmd
}

func aCommand() []string {
	return []string{"sleep", fmt.Sprintf("%d", time.Now().Nanosecond())}
}

func aCustomEnv() map[string]string {
	return map[string]string{"KEY": fmt.Sprintf("%d", time.Now().Nanosecond())}
}

func expectedEnvFor(e map[string]string) []string {
	return toEnvVarsArray(e)
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
