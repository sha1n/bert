package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

type commandExecutor struct {
	pipeStdout bool
	pipeStderr bool
}

// NewCommandExecutor creates a new command executor.
func NewCommandExecutor(pipeStdout bool, pipeStderr bool) api.CommandExecutor {
	return &commandExecutor{
		pipeStdout: pipeStdout,
		pipeStderr: pipeStderr,
	}
}

// Executes a single command in a subprocess based on the specified specs.
func (ce *commandExecutor) Execute(cmdSpec *api.CommandSpec, defaultWorkingDir string, env map[string]string) (exitError error) {
	if cmdSpec == nil {
		return nil
	}

	log.Debugf("Going to execute command %v", cmdSpec.Cmd)

	execCmd := exec.Command(cmdSpec.Cmd[0], cmdSpec.Cmd[1:]...)
	ce.configureCommand(cmdSpec, execCmd, defaultWorkingDir, env)

	exitError = execCmd.Run()

	if exitError != nil {
		log.Errorf("Command '%s' failed. Error: %s", cmdSpec.Cmd, exitError.Error())
	}

	return exitError
}

func (ce *commandExecutor) configureCommand(cmd *api.CommandSpec, execCmd *exec.Cmd, defaultWorkingDir string, env map[string]string) {
	if cmd.WorkingDirectory != "" {
		log.Debugf("Setting command working directory to '%s'", cmd.WorkingDirectory)
		execCmd.Dir = expandPath(cmd.WorkingDirectory)
	} else {
		if defaultWorkingDir != "" {
			log.Debugf("Setting command working directory to '%s'", defaultWorkingDir)
			execCmd.Dir = expandPath(defaultWorkingDir)
		}
	}

	if env != nil {
		cmdEnv := toEnvVarsArray(env)
		log.Debugf("Populating command environment variables '%v'", cmdEnv)
		execCmd.Env = append(execCmd.Env, cmdEnv...)
		execCmd.Env = append(execCmd.Env, os.Environ()...)
	}

	if ce.pipeStdout {
		execCmd.Stdout = log.StandardLogger().Out
	}
	if ce.pipeStdout {
		execCmd.Stderr = log.StandardLogger().Out
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if p, err := os.UserHomeDir(); err == nil {
			return filepath.Join(p, path[1:])
		}
		log.Warnf("Failed to resolve user home for path '%s'", path)
	}

	return path
}

func toEnvVarsArray(env map[string]string) []string {
	var arr []string
	for name, value := range env {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}

	return arr
}
