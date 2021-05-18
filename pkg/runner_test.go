package pkg

import (
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
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

func executeWith(spec *api.BenchmarkSpec) *CmdRecordingExecutor {
	recordingCtx := recordingExecutionContext()

	_ = Execute(spec, recordingCtx)

	return recordingCtx.Executor.(*CmdRecordingExecutor)
}

func assertRecordedCommandWith(t *testing.T, scenario *api.ScenarioSpec) func(expected *api.CommandSpec, actual *RecordedExecutionParams) {
	return func(expected *api.CommandSpec, actual *RecordedExecutionParams) {
		assert.Equal(t, expected, actual.Spec)
		assert.Equal(t, scenario.WorkingDirectory, actual.DefaultWorkingDir)
		assert.Equal(t, scenario.Env, actual.Env)
	}
}

func recordingExecutionContext() *api.ExecutionContext {
	return api.NewExecutionContext(
		NewTracer(),
		&CmdRecordingExecutor{},
	)
}

func aBasicSpecWith(alternate bool, executions int) *api.BenchmarkSpec {
	return &api.BenchmarkSpec{
		Executions: executions,
		Alternate:  alternate,
		Scenarios: []*api.ScenarioSpec{
			{
				Name: "scenario A",
				Command: &api.CommandSpec{
					WorkingDirectory: "/dir-a",
					Cmd:              []string{"cmd", "a"},
				},
			},
			{
				Name: "scenario B",
				Command: &api.CommandSpec{
					WorkingDirectory: "/dir-b",
					Cmd:              []string{"cmd", "b"},
				},
			},
		},
	}
}

func aSpecWithSetupAndTeardownCommands(executions int) *api.BenchmarkSpec {
	return &api.BenchmarkSpec{
		Executions: executions,
		Alternate:  false,
		Scenarios: []*api.ScenarioSpec{
			{
				Name:       "scenario",
				BeforeAll:  &api.CommandSpec{Cmd: []string{"before", "all"}},
				AfterAll:   &api.CommandSpec{Cmd: []string{"after", "all"}},
				BeforeEach: &api.CommandSpec{Cmd: []string{"before", "each"}},
				AfterEach:  &api.CommandSpec{Cmd: []string{"after", "each"}},
				Command: &api.CommandSpec{
					WorkingDirectory: "/home",
					Cmd:              []string{"cmd", "args"},
				},
			},
		},
	}
}

func assertFullScenarioStats(t *testing.T, stats api.Stats) {
	assert.NotNil(t, stats)
	assert.Equal(t, 0.0, stats.ErrorRate())
	assertStatValue(t, stats.Min)
	assertStatValue(t, stats.Max)
	assertStatValue(t, stats.Mean)
	assertStatValue(t, stats.Median)
}

func assertStatValue(t *testing.T, get func() (float64, error)) {
	value, err := get()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, value, 0.0)
}

func failingWriteSummaryFn(t *testing.T) api.WriteReportFn {
	return func(summary api.Summary, config *api.BenchmarkSpec, ctx *api.ReportContext) error {
		t.Fail()
		return nil
	}
}
