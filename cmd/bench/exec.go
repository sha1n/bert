package bench

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Context struct {
	tracer Tracer
}

func NewContext(tracer Tracer) *Context {
	return &Context{
		tracer: tracer,
	}
}

func Execute(b *Benchmark, ctx *Context) TracerSummary {
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
			executeScenario(b.Scenarios[si], ctx)
		}
	}
}

func executeSequencially(b *Benchmark, ctx *Context) {
	for si := range b.Scenarios {
		scenario := b.Scenarios[si]
		for i := 1; i <= b.Executions; i++ {
			executeScenario(scenario, ctx)
		}
	}
}

func executeScenario(scenario *Scenario, ctx *Context) {
	log.Printf("Executing scenario '%s'...\r\n", scenario.Name)
	executeCommand(scenario.Before, scenario.WorkingDirectory, scenario.Env, ctx)

	ctx.tracer.Start(scenario)(
		executeCommand(scenario.Command, scenario.WorkingDirectory, scenario.Env, ctx),
	)

	executeCommand(scenario.After, scenario.WorkingDirectory, scenario.Env, ctx)
}

func executeCommand(cmd *Command, wd string, env map[string]string, ctx *Context) (exitError error) {
	if cmd == nil {
		return nil
	}

	log.Printf("Going to execute command %v", cmd.Cmd)
	benchCmd := exec.Command(cmd.Cmd[0], cmd.Cmd[1:]...)
	if wd != "" {
		log.Printf("Setting custom command working directory '%s'", wd)
		benchCmd.Dir = wd
	}
	if env != nil {
		cmdEnv := toEnvVarsArray(env)
		log.Printf("Appending custom command environment variables '%v'", cmdEnv)
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
