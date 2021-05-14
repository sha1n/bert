package internal

import (
	"github.com/sha1n/benchy/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecuteBenchmarkWithMinimalSpec(t *testing.T) {
	spec := aBasicSpecWith(false, 2)

	execRecordingMock := executeWith(spec)

	assert.Equal(t, 4 /* 2 executions * 2 specs */, len(execRecordingMock.RecordedCommandSeq))
	assertFirstScenarioCommand := assertRecordedCommandWith(t, spec.Scenarios[0])
	assertSecondScenarioCommand := assertRecordedCommandWith(t, spec.Scenarios[1])

	assertFirstScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[0])
	assertFirstScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[1])
	assertSecondScenarioCommand(spec.Scenarios[1].Command, execRecordingMock.RecordedCommandSeq[2])
	assertSecondScenarioCommand(spec.Scenarios[1].Command, execRecordingMock.RecordedCommandSeq[3])
}

func TestExecuteBenchmarkWithMinimalAlternateSpec(t *testing.T) {
	spec := aBasicSpecWith(true, 2)

	execRecordingMock := executeWith(spec)

	assert.Equal(t, 4 /* 2 executions * 2 specs */, len(execRecordingMock.RecordedCommandSeq))
	assertFirstScenarioCommand := assertRecordedCommandWith(t, spec.Scenarios[0])
	assertSecondScenarioCommand := assertRecordedCommandWith(t, spec.Scenarios[1])

	assertFirstScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[0])
	assertSecondScenarioCommand(spec.Scenarios[1].Command, execRecordingMock.RecordedCommandSeq[1])
	assertFirstScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[2])
	assertSecondScenarioCommand(spec.Scenarios[1].Command, execRecordingMock.RecordedCommandSeq[3])
}

func TestExecuteBenchmarkWithSetupAndTeardownSpecs(t *testing.T) {
	spec := aSpecWithSetupAndTeardownCommands(2)

	execRecordingMock := executeWith(spec)

	assert.Equal(t, 8, len(execRecordingMock.RecordedCommandSeq))

	assertScenarioCommand := assertRecordedCommandWith(t, spec.Scenarios[0])

	assertScenarioCommand(spec.Scenarios[0].BeforeAll, execRecordingMock.RecordedCommandSeq[0])

	// Execution #1
	assertScenarioCommand(spec.Scenarios[0].BeforeEach, execRecordingMock.RecordedCommandSeq[1])
	assertScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[2])
	assertScenarioCommand(spec.Scenarios[0].AfterEach, execRecordingMock.RecordedCommandSeq[3])

	// Execution #2
	assertScenarioCommand(spec.Scenarios[0].BeforeEach, execRecordingMock.RecordedCommandSeq[4])
	assertScenarioCommand(spec.Scenarios[0].Command, execRecordingMock.RecordedCommandSeq[5])
	assertScenarioCommand(spec.Scenarios[0].AfterEach, execRecordingMock.RecordedCommandSeq[6])

	assertScenarioCommand(spec.Scenarios[0].AfterAll, execRecordingMock.RecordedCommandSeq[7])
}

func executeWith(spec *BenchmarkSpec) *CmdRecordingExecutor {
	recordingCtx := recordingExecutionContext()

	_ = ExecuteBenchmark(spec, recordingCtx)

	return recordingCtx.executor.(*CmdRecordingExecutor)
}

func assertRecordedCommandWith(t *testing.T, scenario *ScenarioSpec) func(expected *CommandSpec, actual *RecordedExecutionParams) {
	return func(expected *CommandSpec, actual *RecordedExecutionParams) {
		assert.Equal(t, expected, actual.Spec)
		assert.Equal(t, scenario.WorkingDirectory, actual.DefaultWorkingDir)
		assert.Equal(t, scenario.Env, actual.Env)
	}
}

func recordingExecutionContext() *ExecutionContext {
	return NewExecutionContext(
		pkg.NewTracer(),
		&CmdRecordingExecutor{},
	)
}

func aBasicSpecWith(alternate bool, executions int) *BenchmarkSpec {
	return &BenchmarkSpec{
		Executions: executions,
		Alternate:  alternate,
		Scenarios: []*ScenarioSpec{
			{
				Name: "scenario A",
				Command: &CommandSpec{
					WorkingDirectory: "/dir-a",
					Cmd:              []string{"cmd", "a"},
				},
			},
			{
				Name: "scenario B",
				Command: &CommandSpec{
					WorkingDirectory: "/dir-b",
					Cmd:              []string{"cmd", "b"},
				},
			},
		},
	}
}

func aSpecWithSetupAndTeardownCommands(executions int) *BenchmarkSpec {
	return &BenchmarkSpec{
		Executions: executions,
		Alternate:  false,
		Scenarios: []*ScenarioSpec{
			{
				Name:       "scenario",
				BeforeAll:  &CommandSpec{Cmd: []string{"before", "all"}},
				AfterAll:   &CommandSpec{Cmd: []string{"after", "all"}},
				BeforeEach: &CommandSpec{Cmd: []string{"before", "each"}},
				AfterEach:  &CommandSpec{Cmd: []string{"after", "each"}},
				Command: &CommandSpec{
					WorkingDirectory: "/home",
					Cmd:              []string{"cmd", "args"},
				},
			},
		},
	}
}
