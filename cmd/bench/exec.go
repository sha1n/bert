package bench

import (
	"fmt"
	"os"
	"os/exec"
)

type ctx struct {
}

type Context interface {
}

func Execute(b *Benchmark) {
	ctx := &ctx{}

	if b.Alternate {
		executeAlternately(b, ctx)
	} else {
		executeSequencially(b, ctx)
	}

}

func executeAlternately(b *Benchmark, ctx Context) {
	for i := 1; i <= b.Executions; i++ {
		for si := range b.Scenarios {
			executeScenario(b.Scenarios[si], ctx)
		}
	}
}

func executeSequencially(b *Benchmark, ctx Context) {
	for si := range b.Scenarios {
		scenario := b.Scenarios[si]
		for i := 1; i <= b.Executions; i++ {
			executeScenario(scenario, ctx)
		}
	}
}

func executeScenario(scenario *Scenario, ctx Context) {
	fmt.Println(fmt.Sprintf("Executing %s...", scenario.Name))
	executeCommand(scenario.Before, scenario.Env, ctx)

	for ci := range scenario.Script {
		executeCommand(scenario.Script[ci], scenario.Env, ctx)
	}

	executeCommand(scenario.After, scenario.Env, ctx)
}

func executeCommand(cmd *Command, env map[string]string, ctx Context) {
	benchCmd := exec.Command(cmd.Cmd[0], cmd.Cmd[1:]...)
	benchCmd.Stdout = os.Stdout
	benchCmd.Stderr = os.Stderr
	benchCmd.Run()
}
