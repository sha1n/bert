package pkg

import (
	"github.com/sha1n/benchy/api"
)

// Execute executes a benchmark and returns an object that provides access to collected stats.
func Execute(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
	ctx.Listener.OnBenchmarkStart()
	defer ctx.Listener.OnBenchmarkEnd()

	if spec.Alternate {
		executeAlternately(spec, ctx)
	} else {
		executeSequentially(spec, ctx)
	}
}

func executeAlternately(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
	for i := 1; i <= spec.Executions; i++ {
		for si := range spec.Scenarios {
			scenario := spec.Scenarios[si]

			ctx.Listener.OnScenarioStart(scenario.ID())
			if i == 1 {
				executeScenarioSetup(scenario, ctx)
			}
			executeScenarioCommand(scenario, i, spec.Executions, ctx)
			if i == spec.Executions {
				executeScenarioTeardown(scenario, ctx)
			}

			ctx.Listener.OnScenarioEnd(scenario.ID())
		}
	}
}

func executeSequentially(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
	for si := range spec.Scenarios {
		scenario := spec.Scenarios[si]

		for i := 1; i <= spec.Executions; i++ {
			ctx.Listener.OnScenarioStart(scenario.ID())
			if i == 1 {
				executeScenarioSetup(scenario, ctx)
			}

			executeScenarioCommand(scenario, i, spec.Executions, ctx)

			if i == spec.Executions {
				executeScenarioTeardown(scenario, ctx)
			}
			ctx.Listener.OnScenarioEnd(scenario.ID())
		}
	}
}

func executeScenarioSetup(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "running 'beforeAll' command %v...", scenario.BeforeAll.Cmd)
		reportIfError(ctx.Executor.Execute(scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env), scenario.ID(), ctx)
	}
}

func executeScenarioTeardown(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.AfterAll != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "running 'afterAll' command %v...", scenario.AfterAll.Cmd)
		reportIfError(ctx.Executor.Execute(scenario.AfterAll, scenario.WorkingDirectory, scenario.Env), scenario.ID(), ctx)
	}
}

func executeScenarioCommand(scenario api.ScenarioSpec, execIndex int, totalExec int, ctx api.ExecutionContext) {
	ctx.Listener.OnMessagef(scenario.ID(), "run %d of %d", execIndex, totalExec)
	if scenario.BeforeEach != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "running 'beforeEach' command %v", scenario.BeforeEach.Cmd)
		reportIfError(ctx.Executor.Execute(scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env), scenario.ID(), ctx)
	}

	ctx.Listener.OnMessagef(scenario.ID(), "running benchamrk command %v", scenario.Command.Cmd)
	endTrace := ctx.Tracer.Start(scenario)
	err := ctx.Executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env)
	reportIfError(err, scenario.ID(), ctx)
	endTrace(err)

	if scenario.AfterEach != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "running 'afterEach' command %v", scenario.AfterEach.Cmd)
		reportIfError(ctx.Executor.Execute(scenario.AfterEach, scenario.WorkingDirectory, scenario.Env), scenario.ID(), ctx)
	}
}

func reportIfError(err error, id api.ID, ctx api.ExecutionContext) {
	if err != nil {
		ctx.Listener.OnError(id, err)
	}
}
