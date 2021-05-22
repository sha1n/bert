package cli

import (
	"fmt"
	"strings"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateConfig(cmd *cobra.Command, args []string) {
	printHints()
	outFile := ResolveOutputFileArg(cmd, ArgNameOutputFile)

	spec := &api.BenchmarkSpec{
		Executions: int(RequestUint("number of executions", true)),
		Alternate:  RequestBool("alternate executions", false),
		Scenarios:  RequestScenarios(),
	}

	fmt.Print(`

Writing your configuration...

`)

	if err := pkg.SaveSpec(spec, outFile.Name()); err != nil {
		log.Error(err)
		log.Exit(1)
	}
}

func RequestScenarios() []*api.ScenarioSpec {
	specs := []*api.ScenarioSpec{}

	for {
		specs = append(specs, RequestScenario())
		if !QuestionYN("add another scenario?") {
			break
		}
	}

	return specs
}

func RequestCommand(description string, required bool) *api.CommandSpec {
	requestCommand := func() *api.CommandSpec {
		return &api.CommandSpec{
			WorkingDirectory: RequestExistingDirectory("working directory", false),
			Cmd:              RequestCommandLine("command line", true),
		}
	}

	if required {
		_, _ = printfBold("%s:\r\n", description)
		return requestCommand()
	}
	if QuestionYN(fmt.Sprintf("%s?", description)) {
		return requestCommand()
	}

	return nil
}

func RequestEnvVars() map[string]string {
	var envVars map[string]string

	if QuestionYN("define custom env vars?") {
		envVars = map[string]string{}
		for {
			kv := RequestString("K=v", false)
			if kv != "" {
				kvSlice := strings.Split(kv, "=")
				envVars[kvSlice[0]] = kvSlice[1]
			} else {
				break
			}
		}
	}

	return envVars
}

func RequestScenario() *api.ScenarioSpec {
	return &api.ScenarioSpec{
		Name:             RequestString("scenario name", true),
		WorkingDirectory: RequestExistingDirectory("working directory", false),
		Env:              RequestEnvVars(),
		BeforeAll:        RequestCommand("add setup command", false),
		AfterAll:         RequestCommand("add teardown command", false),
		BeforeEach:       RequestCommand("add before each command", false),
		AfterEach:        RequestCommand("add after each command", false),
		Command:          RequestCommand("benchmarked command", true),
	}
}

func printHints() {
	fmt.Printf(`
--------------------------------
 BENCHMARK CONFIGURATION HELPER
--------------------------------

This tool is going to help you go through a benchmark configuration definition.

%s annotates required input 
%s annotates optional input

more here: https://github.com/sha1n/benchy/blob/master/docs/configuration.md

--------------------------------

`,
		sprintBold(sprintRed("*")),
		sprintBold(sprintGreen("?")),
	)
}
