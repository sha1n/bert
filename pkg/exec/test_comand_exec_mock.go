package exec

import (
	"context"

	"github.com/sha1n/bert/api"
)

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
	ctx context.Context,
	cmdSpec *api.CommandSpec,
	defaultWorkingDir string,
	env map[string]string,
) api.ExecCommandFn {

	ce.RecordedCommandSeq = append(ce.RecordedCommandSeq, &RecordedExecutionParams{
		Spec:              cmdSpec,
		DefaultWorkingDir: defaultWorkingDir,
		Env:               env,
	})

	return func() (*api.ExecutionInfo, error) { return &api.ExecutionInfo{}, nil }
}
