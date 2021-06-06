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

			if i == 1 {
				executeScenarioSetup(scenario, ctx)
			}
			executeScenarioCommand(scenario, i, spec.Executions, ctx)
			if i == spec.Executions {
				executeScenarioTeardown(scenario, ctx)
			}
		}
	}
}

func executeSequentially(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
	for si := range spec.Scenarios {
		scenario := spec.Scenarios[si]

		executeScenarioSetup(scenario, ctx)
		for i := 1; i <= spec.Executions; i++ {
			executeScenarioCommand(scenario, i, spec.Executions, ctx)
		}
		executeScenarioTeardown(scenario, ctx)
	}
}

func executeScenarioSetup(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		log.Debugf("Running setup for scenario '%s'...", scenario.Name)
		logError(ctx.Executor.Execute(scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env))
	}
}

func executeScenarioTeardown(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.AfterAll != nil {
		log.Debugf("Running teardown for scenario '%s'...", scenario.Name)
		logError(ctx.Executor.Execute(scenario.AfterAll, scenario.WorkingDirectory, scenario.Env))
	}
}

func executeScenarioCommand(scenario api.ScenarioSpec, execIndex int, totalExec int, ctx api.ExecutionContext) {
	log.Infof("Executing scenario '%s'... (%d/%d)", scenario.Name, execIndex, totalExec)
	if scenario.BeforeEach != nil {
		log.Debugf("Executing 'before' command %v", scenario.BeforeEach.Cmd)
		logError(ctx.Executor.Execute(scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env))
	}

	ctx.Tracer.Start(scenario)(
		ctx.Executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env),
	)

	if scenario.AfterEach != nil {
		log.Debugf("Executing 'after' command %v", scenario.AfterEach.Cmd)
		logError(ctx.Executor.Execute(scenario.AfterEach, scenario.WorkingDirectory, scenario.Env))
	}
}

func logError(err error) {
	if err != nil {
		log.Error(err)
	}
}
