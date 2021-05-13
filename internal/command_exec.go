package internal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

type CommandExecutor interface {
	Execute(cmd *CommandSpec, defaultWorkingDir string, env map[string]string, ctx *ExecutionContext) error
}

type commandExecutor struct {
	pipeStdout bool
	pipeStderr bool
}

func NewCommandExecutor(pipeStdout bool, pipeStderr bool) CommandExecutor {
	return &commandExecutor{
		pipeStdout: pipeStdout,
		pipeStderr: pipeStderr,
	}
}

func NewDefaultCommandExecutor() CommandExecutor {
	return NewCommandExecutor(true, true)
}

func (ce *commandExecutor) Execute(cmdSpec *CommandSpec, defaultWorkingDir string, env map[string]string, ctx *ExecutionContext) (exitError error) {
	if cmdSpec == nil {
		return nil
	}

	log.Debugf("Going to execute command %v", cmdSpec.Cmd)

	execCmd := exec.Command(cmdSpec.Cmd[0], cmdSpec.Cmd[1:]...)
	ce.configureCommand(cmdSpec, execCmd, defaultWorkingDir, env, ctx)

	exitError = execCmd.Run()

	if exitError != nil {
		log.Errorf("Command '%s' failed. Error: %s", cmdSpec.Cmd, exitError.Error())
	}

	return exitError
}

func (ce *commandExecutor) configureCommand(cmd *CommandSpec, execCmd *exec.Cmd, defaultWorkingDir string, env map[string]string, ctx *ExecutionContext) {
	if cmd.WorkingDirectory != "" {
		log.Debugf("Setting command working directory to '%s'", cmd.WorkingDirectory)
		execCmd.Dir = cmd.WorkingDirectory
	} else {
		if defaultWorkingDir != "" {
			log.Debugf("Setting command working directory to '%s'", defaultWorkingDir)
			execCmd.Dir = defaultWorkingDir
		}
	}

	if env != nil {
		cmdEnv := toEnvVarsArray(env)
		log.Debugf("Populating command environment variables '%v'", cmdEnv)
		execCmd.Env = append(execCmd.Env, cmdEnv...)
	}

	if ce.pipeStdout {
		execCmd.Stdout = os.Stdout
	}
	if ce.pipeStdout {
		execCmd.Stderr = os.Stderr
	}
}

func toEnvVarsArray(env map[string]string) []string {
	var arr []string
	for name, value := range env {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}

	return arr
}
