package api

import "time"

// ExecutionInfo information about an executed command
type ExecutionInfo struct {
	UserTime      time.Duration
	SystemTime    time.Duration
	PerceivedTime time.Duration
	ExitCode      int
}

// ExecCommandFn executes a command and returns execution information or an error
// if the command failed to execute or returned a non-zero exit code.
type ExecCommandFn = func() (*ExecutionInfo, error)

// CommandExecutor is an abstraction for commands executed as subprocesses.
type CommandExecutor interface {
	ExecuteFn(cmd *CommandSpec, defaultWorkingDir string, env map[string]string) ExecCommandFn
}

// ExecutionContext provides access to benchmark global resources
type ExecutionContext struct {
	Executor CommandExecutor
	Tracer   Tracer
	Listener
}

// NewExecutionContext creates a new ExecutionContext.
func NewExecutionContext(tracer Tracer, executor CommandExecutor, listener Listener) ExecutionContext {
	return ExecutionContext{
		Executor: executor,
		Tracer:   tracer,
		Listener: listener,
	}
}
