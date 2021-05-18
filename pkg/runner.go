package pkg

import (
	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// Run runs a benchmark.
func Run(specFilePath string, executioinCtx *api.ExecutionContext, write api.WriteReportFn) (err error) {
	log.Info("Starting benchy...")

	var spec *api.BenchmarkSpec
	if spec, err = loadSpec(specFilePath); err == nil {
		summary := Execute(spec, executioinCtx)
		err = write(summary, spec)
	}

	return err
}

func loadSpec(filePath string) (*api.BenchmarkSpec, error) {
	log.Infof("Loading benchmark specs from '%s'...", filePath)

	return LoadSpec(filePath)
}

// Execute executes a benchmark and returns an object that provides access to collected stats.
func Execute(spec *api.BenchmarkSpec, ctx *api.ExecutionContext) api.Summary {
	if spec.Alternate {
		executeAlternately(spec, ctx)
	} else {
		executeSequencially(spec, ctx)
	}

	return ctx.Tracer.Summary()
}

func executeAlternately(b *api.BenchmarkSpec, ctx *api.ExecutionContext) {
	for i := 1; i <= b.Executions; i++ {
		for si := range b.Scenarios {
			scenario := b.Scenarios[si]

			if i == 1 {
				executeScenarioSetup(scenario, ctx)
			}
			executeScenarioCommand(scenario, i, b.Executions, ctx)
			if i == b.Executions {
				executeScenarioTeardown(scenario, ctx)
			}
		}
	}
}

func executeSequencially(b *api.BenchmarkSpec, ctx *api.ExecutionContext) {
	for si := range b.Scenarios {
		scenario := b.Scenarios[si]

		executeScenarioSetup(scenario, ctx)
		for i := 1; i <= b.Executions; i++ {
			executeScenarioCommand(scenario, i, b.Executions, ctx)
		}
		executeScenarioTeardown(scenario, ctx)
	}
}

func executeScenarioSetup(scenario *api.ScenarioSpec, ctx *api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		log.Debugf("Running setup for scenario '%s'...", scenario.Name)
		ctx.Executor.Execute(scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env)
	}
}

func executeScenarioTeardown(scenario *api.ScenarioSpec, ctx *api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		log.Debugf("Running teardown for scenario '%s'...", scenario.Name)
		ctx.Executor.Execute(scenario.AfterAll, scenario.WorkingDirectory, scenario.Env)
	}
}

func executeScenarioCommand(scenario *api.ScenarioSpec, execIndex int, totalExec int, ctx *api.ExecutionContext) {
	log.Infof("Executing scenario '%s'... (%d/%d)", scenario.Name, execIndex, totalExec)
	if scenario.BeforeEach != nil {
		log.Debugf("Executing 'before' command %v", scenario.BeforeEach.Cmd)
		ctx.Executor.Execute(scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env)
	}

	ctx.Tracer.Start(scenario)(
		ctx.Executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env),
	)

	if scenario.AfterEach != nil {
		log.Debugf("Executing 'after' command %v", scenario.AfterEach.Cmd)
		ctx.Executor.Execute(scenario.AfterEach, scenario.WorkingDirectory, scenario.Env)
	}
}
