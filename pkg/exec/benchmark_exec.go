package exec

import (
	"context"

	"github.com/sha1n/bert/api"
)

// Execute executes a benchmark and returns an object that provides access to collected stats.
func Execute(ctx context.Context, spec api.BenchmarkSpec, execCtx api.ExecutionContext) {
	execCtx.OnBenchmarkStart()
	defer execCtx.OnBenchmarkEnd()

	if spec.Alternate {
		executeAlternately(ctx, spec, execCtx)
	} else {
		executeSequentially(ctx, spec, execCtx)
	}
}

func executeAlternately(ctx context.Context, spec api.BenchmarkSpec, execCtx api.ExecutionContext) {
	for i := 1; i <= spec.Executions; i++ {
		for si := range spec.Scenarios {
			if ctx.Err() != nil {
				return
			}

			scenario := spec.Scenarios[si]

			execCtx.OnScenarioStart(scenario.ID())
			if i == 1 {
				executeScenarioSetup(ctx, scenario, execCtx)
			}
			executeScenarioCommand(ctx, scenario, i, spec.Executions, execCtx)
			if i == spec.Executions {
				executeScenarioTeardown(ctx, scenario, execCtx)
			}

			execCtx.OnScenarioEnd(scenario.ID())
		}
	}
}

func executeSequentially(ctx context.Context, spec api.BenchmarkSpec, execCtx api.ExecutionContext) {
	for si := range spec.Scenarios {
		scenario := spec.Scenarios[si]

		for i := 1; i <= spec.Executions; i++ {
			if ctx.Err() != nil {
				return
			}

			execCtx.OnScenarioStart(scenario.ID())
			if i == 1 {
				executeScenarioSetup(ctx, scenario, execCtx)
			}

			executeScenarioCommand(ctx, scenario, i, spec.Executions, execCtx)

			if i == spec.Executions {
				executeScenarioTeardown(ctx, scenario, execCtx)
			}
			execCtx.OnScenarioEnd(scenario.ID())
		}
	}
}

func executeScenarioSetup(ctx context.Context, scenario api.ScenarioSpec, execCtx api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		execCtx.OnMessagef(scenario.ID(), "running 'beforeAll' command %v...", scenario.BeforeAll.Cmd)
		reportIfExecError(execCtx.Executor.ExecuteFn(ctx, scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env), scenario.ID(), execCtx)
	}
}

func executeScenarioTeardown(ctx context.Context, scenario api.ScenarioSpec, execCtx api.ExecutionContext) {
	if scenario.AfterAll != nil {
		execCtx.OnMessagef(scenario.ID(), "running 'afterAll' command %v...", scenario.AfterAll.Cmd)
		reportIfExecError(execCtx.Executor.ExecuteFn(ctx, scenario.AfterAll, scenario.WorkingDirectory, scenario.Env), scenario.ID(), execCtx)
	}
}

func executeScenarioCommand(ctx context.Context, scenario api.ScenarioSpec, execIndex int, totalExec int, execCtx api.ExecutionContext) {
	execCtx.OnMessagef(scenario.ID(), "run %d of %d", execIndex, totalExec)
	if scenario.BeforeEach != nil {
		execCtx.OnMessagef(scenario.ID(), "running 'beforeEach' command %v", scenario.BeforeEach.Cmd)
		reportIfExecError(execCtx.Executor.ExecuteFn(ctx, scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env), scenario.ID(), execCtx)
	}

	execCtx.OnMessagef(scenario.ID(), "running benchmark command %v", scenario.Command.Cmd)
	executeFn := execCtx.Executor.ExecuteFn(ctx, scenario.Command, scenario.WorkingDirectory, scenario.Env)

	endTrace := execCtx.Tracer.Start(scenario)
	info, err := executeFn()
	endTrace(info, err)

	reportIfError(err, scenario.ID(), execCtx)

	if scenario.AfterEach != nil {
		execCtx.OnMessagef(scenario.ID(), "running 'afterEach' command %v", scenario.AfterEach.Cmd)
		reportIfExecError(execCtx.Executor.ExecuteFn(ctx, scenario.AfterEach, scenario.WorkingDirectory, scenario.Env), scenario.ID(), execCtx)
	}
}

func reportIfError(err error, id api.ID, ctx api.ExecutionContext) {
	if err != nil {
		ctx.OnError(id, err)
	}
}

func reportIfExecError(exec api.ExecCommandFn, id api.ID, ctx api.ExecutionContext) {
	if _, err := exec(); err != nil {
		ctx.OnError(id, err)
	}
}
