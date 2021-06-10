package pkg

import (
	"fmt"
	"os"
	"os/exec"

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
func (ce *commandExecutor) ExecuteFn(cmdSpec *api.CommandSpec, defaultWorkingDir string, env map[string]string) func() error {
	log.Debugf("Going to execute command %v", cmdSpec.Cmd)

	execCmd := exec.Command(cmdSpec.Cmd[0], cmdSpec.Cmd[1:]...)
	ce.configureCommand(cmdSpec, execCmd, defaultWorkingDir, env)

	cancel, _ := RegisterInterruptGuard(onShutdownSignalFn(execCmd))

	return func() error {
		defer cancel()
		return execCmd.Run()
	}
}

func (ce *commandExecutor) configureCommand(cmd *api.CommandSpec, execCmd *exec.Cmd, defaultWorkingDir string, env map[string]string) {
	if cmd.WorkingDirectory != "" {
		log.Debugf("Setting command working directory to '%s'", cmd.WorkingDirectory)
		execCmd.Dir = ExpandUserPath(cmd.WorkingDirectory)
	} else {
		if defaultWorkingDir != "" {
			log.Debugf("Setting command working directory to '%s'", defaultWorkingDir)
			execCmd.Dir = ExpandUserPath(defaultWorkingDir)
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
	if ce.pipeStderr {
		execCmd.Stderr = log.StandardLogger().Out
	}
}

func toEnvVarsArray(env map[string]string) []string {
	var arr []string
	for name, value := range env {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}

	return arr
}

func onShutdownSignalFn(execCmd *exec.Cmd) func(os.Signal) {
	return func(sig os.Signal) {
		if sig == os.Interrupt {
			log.Debugf("Got %s signal. Forwarding to %s...", sig, execCmd.Args[0])
			execCmd.Process.Signal(sig)

			os.Exit(1)
		}
	}
}
