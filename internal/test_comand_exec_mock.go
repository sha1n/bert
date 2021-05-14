package internal

// CmdRecordingExecutor ...
type CmdRecordingExecutor struct {
	RecordedCommandSeq []*RecordedExecutionParams
}

// RecordedExecutionParams ...
type RecordedExecutionParams struct {
	Spec              *CommandSpec
	DefaultWorkingDir string
	Env               map[string]string
	Ctx               *ExecutionContext
}

// Execute records execution parameters and stores them in order
func (ce *CmdRecordingExecutor) Execute(
	cmdSpec *CommandSpec,
	defaultWorkingDir string,
	env map[string]string,
	ctx *ExecutionContext) (exitError error) {

	ce.RecordedCommandSeq = append(ce.RecordedCommandSeq, &RecordedExecutionParams{
		Spec:              cmdSpec,
		DefaultWorkingDir: defaultWorkingDir,
		Env:               env,
		Ctx:               ctx,
	})

	return nil
}
