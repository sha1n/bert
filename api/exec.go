package api

// CommandExecutor is an abstraction for commands executed as subprocesses.
type CommandExecutor interface {
	Execute(cmd *CommandSpec, defaultWorkingDir string, env map[string]string) error
}

// ExecutionContext provides access to benchmark global resources
type ExecutionContext struct {
	Executor CommandExecutor
	Tracer   Tracer
}

// NewExecutionContext creates a new ExecutionContext.
func NewExecutionContext(tracer Tracer, executor CommandExecutor) ExecutionContext {
	return ExecutionContext{
		Executor: executor,
		Tracer:   tracer,
	}
}
