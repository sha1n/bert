package exec

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg/osutil"
)

type commandExecutor struct {
	pipeStdout   bool
	pipeStderr   bool
	outputWriter io.Writer
}

// NewCommandExecutor creates a new command executor.
func NewCommandExecutor(pipeStdout bool, pipeStderr bool, outputWriter io.Writer) api.CommandExecutor {
	return &commandExecutor{
		pipeStdout:   pipeStdout,
		pipeStderr:   pipeStderr,
		outputWriter: outputWriter,
	}
}

// Executes a single command in a subprocess based on the specified specs.
func (ce *commandExecutor) ExecuteFn(ctx context.Context, cmdSpec *api.CommandSpec, defaultWorkingDir string, env map[string]string) api.ExecCommandFn {
	slog.Debug(fmt.Sprintf("Going to execute command %v", cmdSpec.Cmd))

	execCmd := exec.CommandContext(ctx, cmdSpec.Cmd[0], cmdSpec.Cmd[1:]...)
	ce.configureCommand(cmdSpec, execCmd, defaultWorkingDir, env)

	return func() (execInfo *api.ExecutionInfo, err error) {
		startTime := time.Now()
		err = execCmd.Run()
		perceivedTime := time.Since(startTime)
		state := execCmd.ProcessState
		if state != nil && state.Exited() {
			execInfo = &api.ExecutionInfo{
				ExitCode:      state.ExitCode(),
				UserTime:      state.UserTime(),
				SystemTime:    state.SystemTime(),
				PerceivedTime: perceivedTime,
			}

		}

		return
	}
}

func (ce *commandExecutor) configureCommand(cmd *api.CommandSpec, execCmd *exec.Cmd, defaultWorkingDir string, env map[string]string) {
	if cmd.WorkingDirectory != "" {
		slog.Debug(fmt.Sprintf("Setting command working directory to '%s'", cmd.WorkingDirectory))
		execCmd.Dir = osutil.ExpandUserPath(cmd.WorkingDirectory)
	} else {
		if defaultWorkingDir != "" {
			slog.Debug(fmt.Sprintf("Setting command working directory to '%s'", defaultWorkingDir))
			execCmd.Dir = osutil.ExpandUserPath(defaultWorkingDir)
		}
	}

	if env != nil {
		cmdEnv := toEnvVarsArray(env)
		slog.Debug(fmt.Sprintf("Populating command environment variables '%v'", cmdEnv))
		execCmd.Env = append(execCmd.Env, os.Environ()...)
		execCmd.Env = append(execCmd.Env, cmdEnv...)
	}

	if ce.pipeStdout {
		execCmd.Stdout = ce.outputWriter
	}
	if ce.pipeStderr {
		execCmd.Stderr = ce.outputWriter
	}
}

func toEnvVarsArray(env map[string]string) []string {
	var arr []string
	for name, value := range env {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}

	return arr
}
