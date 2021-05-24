package cli

import (
	"fmt"
	"strings"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateConfigCommand creates the 'config' sub command
func CreateConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Long:  `Interactively walks through a benchmark configuration creation process`,
		Short: `Interactively creates a benchmark config`,
		Run:   CreateConfig,
	}

	cmd.Flags().StringP(ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)

	_ = cmd.MarkFlagFilename(ArgNameOutputFile, "yml", "yaml")

	return cmd
}

func CreateConfig(cmd *cobra.Command, args []string) {
	printHints()
	writeCloser := ResolveOutputArg(cmd, ArgNameOutputFile)
	defer writeCloser.Close()

	spec := &api.BenchmarkSpec{
		Executions: int(RequestUint("number of executions", true)),
		Alternate:  RequestOptionalBool("alternate executions", false),
		Scenarios:  RequestScenarios(),
	}

	fmt.Print(`

Writing your configuration...

`)

	if err := pkg.SaveSpec(spec, writeCloser); err != nil {
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
			WorkingDirectory: RequestOptionalExistingDirectory("working directory", "inherits scenario"),
			Cmd:              RequestCommandLine("command line", true),
		}
	}

	if required {
		_, _ = fmt.Printf("%s:\r\n", description)
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
		WorkingDirectory: RequestOptionalExistingDirectory("working directory", "inherits benchy's"),
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
