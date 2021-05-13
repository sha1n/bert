package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sha1n/benchy/pkg"
)

type Context struct {
	tracer pkg.Tracer
}

func NewContext(tracer pkg.Tracer) *Context {
	return &Context{
		tracer: tracer,
	}
}

func Execute(b *Benchmark, ctx *Context) pkg.TracerSummary {
	if b.Alternate {
		executeAlternately(b, ctx)
	} else {
		executeSequencially(b, ctx)
	}

	return ctx.tracer.Summary()
}

func executeAlternately(b *Benchmark, ctx *Context) {
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

func executeSequencially(b *Benchmark, ctx *Context) {
	for si := range b.Scenarios {
		scenario := b.Scenarios[si]

		executeScenarioSetup(scenario, ctx)
		for i := 1; i <= b.Executions; i++ {
			executeScenarioCommand(scenario, ctx)
		}
		executeScenarioTeardown(scenario, ctx)
	}
}

func executeScenarioSetup(scenario *Scenario, ctx *Context) {
	log.Printf("Running setup for scenario '%s'...\r\n", scenario.Name)
	executeCommand(scenario.Setup, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeScenarioTeardown(scenario *Scenario, ctx *Context) {
	log.Printf("Running teardown for scenario '%s'...\r\n", scenario.Name)
	executeCommand(scenario.Teardown, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeScenarioCommand(scenario *Scenario, ctx *Context) {
	log.Printf("Executing scenario '%s'...\r\n", scenario.Name)
	executeCommand(scenario.BeforeCommand, scenario.WorkingDirectory, scenario.Env, ctx)

	ctx.tracer.Start(scenario)(
		executeCommand(scenario.Command, scenario.WorkingDirectory, scenario.Env, ctx),
	)

	executeCommand(scenario.AfterCommand, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeCommand(cmd *Command, defaultWorkingDir string, env map[string]string, ctx *Context) (exitError error) {
	if cmd == nil {
		return nil
	}

	log.Printf("Going to execute command %v", cmd.Cmd)

	benchCmd := exec.Command(cmd.Cmd[0], cmd.Cmd[1:]...)

	if cmd.WorkingDirectory != "" {
		log.Printf("Setting command working directory to '%s'", cmd.WorkingDirectory)
		benchCmd.Dir = cmd.WorkingDirectory
	} else {
		if defaultWorkingDir != "" {
			log.Printf("Setting command working directory to '%s'", defaultWorkingDir)
			benchCmd.Dir = defaultWorkingDir
		}
	}

	if env != nil {
		cmdEnv := toEnvVarsArray(env)
		log.Printf("Populating command environment variables '%v'", cmdEnv)
		benchCmd.Env = append(benchCmd.Env, cmdEnv...)
	}

	benchCmd.Stdout = os.Stdout
	benchCmd.Stderr = os.Stderr

	exitError = benchCmd.Run()

	if exitError != nil {
		log.Printf("[ERROR] Command '%s' failed. Error: %s", cmd.Cmd, exitError.Error())
	}

	return exitError
}

func toEnvVarsArray(env map[string]string) []string {
	var arr []string
	for name, value := range env {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}

	return arr
}
