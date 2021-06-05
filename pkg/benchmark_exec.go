package pkg

import (
	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// Execute executes a benchmark and returns an object that provides access to collected stats.
func Execute(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
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
		
		ctx.Listener.OnScenarioStart(scenario.ID())
		
		executeScenarioSetup(scenario, ctx)
		for i := 1; i <= spec.Executions; i++ {
			executeScenarioCommand(scenario, i, spec.Executions, ctx)
		}
		executeScenarioTeardown(scenario, ctx)
		
		ctx.Listener.OnScenarioEnd(scenario.ID())
	}
}

func executeScenarioSetup(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		ctx.Listener.OnMessage(scenario.ID(), "Running setup...")
		logError(scenario.ID(), ctx.Executor.Execute(scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env), ctx)
	}
}

func executeScenarioTeardown(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.AfterAll != nil {
		ctx.Listener.OnMessage(scenario.ID(), "Running teardown...")
		logError(scenario.ID(), ctx.Executor.Execute(scenario.AfterAll, scenario.WorkingDirectory, scenario.Env), ctx)
	}
}

func executeScenarioCommand(scenario api.ScenarioSpec, execIndex int, totalExec int, ctx api.ExecutionContext) {
	ctx.Listener.OnMessagef(scenario.ID(), "Running benchmarked command... (%d/%d)", execIndex, totalExec)
	if scenario.BeforeEach != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "Running 'before' command %v", scenario.BeforeEach.Cmd)
		logError(scenario.ID(), ctx.Executor.Execute(scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env), ctx)
	}

	ctx.Tracer.Start(scenario)(
		ctx.Executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env),
	)

	if scenario.AfterEach != nil {
		ctx.Listener.OnMessagef(scenario.ID(), "Running 'after' command %v", scenario.AfterEach.Cmd)
		logError(scenario.ID(), ctx.Executor.Execute(scenario.AfterEach, scenario.WorkingDirectory, scenario.Env), ctx)
	}
}

func logError(id api.ID, err error, ctx api.ExecutionContext) {
	if err != nil {
		ctx.Listener.OnError(id, err)
		log.Error(err)
	}
}
