package pkg

import (
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
)

// Execute executes a benchmark and returns an object that provides access to collected stats.
func Execute(spec api.BenchmarkSpec, ctx api.ExecutionContext) {
	matrix := termite.NewMatrix(ctx.StdoutWriter, time.Millisecond*100)
	rows := matrix.NewRange(len(spec.Scenarios) + 4)

	println()

	// rows[0].WriteString("Executing Benchmark")
	ctx.StdoutWriter = rows[len(spec.Scenarios) + 3]

	tickers := make([]func() bool, len(spec.Scenarios))
	termWidth, _, _ := termite.GetTerminalDimensions()
	for i := 2; i < len(spec.Scenarios) + 2; i++ {
		bar := termite.NewProgressBar(rows[i], spec.Executions, termWidth, 50, &progressBarFormatter{})
		tick, cancel, _ := bar.Start()
		tickers[i - 2] = tick
		defer cancel()
	}

	cancel := matrix.Start()
	defer cancel()

	log.Info("Executing benchmark scenarios...")
	if spec.Alternate {
		executeAlternately(spec, ctx, tickers)
	} else {
		executeSequentially(spec, ctx, tickers)
	}
	log.Info("Finished executing all benchmark scenarios")
}

func executeAlternately(spec api.BenchmarkSpec, ctx api.ExecutionContext, bars []func() bool) {
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

			bars[si]()
		}
	}
}

func executeSequentially(spec api.BenchmarkSpec, ctx api.ExecutionContext, bars []func() bool) {
	for si := range spec.Scenarios {
		scenario := spec.Scenarios[si]
		bars[si]()
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
		logError(ctx.Executor.Execute(scenario.BeforeAll, scenario.WorkingDirectory, scenario.Env, ctx))
	}
}

func executeScenarioTeardown(scenario api.ScenarioSpec, ctx api.ExecutionContext) {
	if scenario.BeforeAll != nil {
		log.Debugf("Running teardown for scenario '%s'...", scenario.Name)
		logError(ctx.Executor.Execute(scenario.AfterAll, scenario.WorkingDirectory, scenario.Env, ctx))
	}
}

func executeScenarioCommand(scenario api.ScenarioSpec, execIndex int, totalExec int, ctx api.ExecutionContext) {
	log.Infof("Executing scenario '%s'... (%d/%d)", scenario.Name, execIndex, totalExec)
	if scenario.BeforeEach != nil {
		log.Debugf("Executing 'before' command %v", scenario.BeforeEach.Cmd)
		logError(ctx.Executor.Execute(scenario.BeforeEach, scenario.WorkingDirectory, scenario.Env, ctx))
	}

	ctx.Tracer.Start(scenario)(
		ctx.Executor.Execute(scenario.Command, scenario.WorkingDirectory, scenario.Env, ctx),
	)

	if scenario.AfterEach != nil {
		log.Debugf("Executing 'after' command %v", scenario.AfterEach.Cmd)
		logError(ctx.Executor.Execute(scenario.AfterEach, scenario.WorkingDirectory, scenario.Env, ctx))
	}
}

func logError(err error) {
	if err != nil {
		log.Error(err)
	}
}

type progressBarFormatter struct {}

// FormatLeftBorder returns the left border char
func (f *progressBarFormatter) FormatLeftBorder() string {
	return cyan.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatRightBorder returns the right border char
func (f *progressBarFormatter) FormatRightBorder() string {
	return cyan.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatFill returns the fill char
func (f *progressBarFormatter) FormatFill() string {
	return cyan.Sprintf("%c", termite.DefaultProgressBarFill)
}