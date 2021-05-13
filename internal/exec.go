package internal

import (
	"log"

	"github.com/sha1n/benchy/pkg"
)

type ExecutionContext struct {
	executor CommandExecutor
	tracer   pkg.Tracer
}

func NewExecutionContext(tracer pkg.Tracer, executor CommandExecutor) *ExecutionContext {
	return &ExecutionContext{
		executor: executor,
		tracer:   tracer,
	}
}

func ExecuteBenchmark(spec *BenchmarkSpec, ctx *ExecutionContext) pkg.TracerSummary {
	if spec.Alternate {
		executeAlternately(spec, ctx)
	} else {
		executeSequencially(spec, ctx)
	}

	return ctx.tracer.Summary()
}

func executeAlternately(b *BenchmarkSpec, ctx *ExecutionContext) {
	for i := 1; i <= b.Executions; i++ {
		for si := range b.Scenarios {
			scenario := b.Scenarios[si]

			if i == 1 {
				executeScenarioSetup(scenario, ctx)
			}
			executeScenarioCommand(scenario, ctx)
			if i == b.Executions {
				executeScenarioTeardown(scenario, ctx)
			}
		}
	}
}

func executeSequencially(b *BenchmarkSpec, ctx *ExecutionContext) {
	for si := range b.Scenarios {
		scenario := b.Scenarios[si]

		executeScenarioSetup(scenario, ctx)
		for i := 1; i <= b.Executions; i++ {
			executeScenarioCommand(scenario, ctx)
		}
		executeScenarioTeardown(scenario, ctx)
	}
}

func executeScenarioSetup(scenario *ScenarioSpec, ctx *ExecutionContext) {
	log.Printf("Running setup for scenario '%s'...\r\n", scenario.Name)
	ctx.executor.Execute(scenario.Setup, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeScenarioTeardown(scenario *ScenarioSpec, ctx *ExecutionContext) {
	log.Printf("Running teardown for scenario '%s'...\r\n", scenario.Name)
	ctx.executor.Execute(scenario.Teardown, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeScenarioCommand(scenario *ScenarioSpec, ctx *ExecutionContext) {
	log.Printf("Executing scenario '%s'...\r\n", scenario.Name)
	ctx.executor.Execute(scenario.BeforeCommand, scenario.WorkingDirectory, scenario.Env, ctx)

	ctx.tracer.Start(scenario)(
		ctx.executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env, ctx),
	)

	ctx.executor.Execute(scenario.AfterCommand, scenario.WorkingDirectory, scenario.Env, ctx)
}
