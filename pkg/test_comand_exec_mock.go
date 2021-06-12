package pkg

import "github.com/sha1n/benchy/api"

// CmdRecordingExecutor ...
type CmdRecordingExecutor struct {
	RecordedCommandSeq []*RecordedExecutionParams
}

// RecordedExecutionParams ...
type RecordedExecutionParams struct {
	Spec              *api.CommandSpec
	DefaultWorkingDir string
	Env               map[string]string
}

// ExecuteFn records execution parameters and stores them in order
func (ce *CmdRecordingExecutor) ExecuteFn(
	cmdSpec *api.CommandSpec,
	defaultWorkingDir string,
	env map[string]string,
) func() error {

	ce.RecordedCommandSeq = append(ce.RecordedCommandSeq, &RecordedExecutionParams{
		Spec:              cmdSpec,
		DefaultWorkingDir: defaultWorkingDir,
		Env:               env,
	})

	return func() error { return nil }
}
