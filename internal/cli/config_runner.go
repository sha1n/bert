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

// CreateConfig runs the config tool
func CreateConfig(cmd *cobra.Command, args []string) {
	printHints()
	writeCloser := ResolveOutputArg(cmd, ArgNameOutputFile)
	defer writeCloser.Close()

	spec := &api.BenchmarkSpec{
		Executions: int(requestUint("number of executions", true)),
		Alternate:  requestOptionalBool("alternate executions", false),
		Scenarios:  requestScenarios(),
	}

	fmt.Print(`

Writing your configuration...

`)

	if err := pkg.SaveSpec(spec, writeCloser); err != nil {
		log.Error(err)
		log.Exit(1)
	}
}

func requestScenarios() []*api.ScenarioSpec {
	specs := []*api.ScenarioSpec{}

	for {
		specs = append(specs, requestScenario())
		if !questionYN("add another scenario?") {
			break
		}
	}

	return specs
}

func requestCommand(description string, required bool) *api.CommandSpec {
	requestCommand := func() *api.CommandSpec {
		return &api.CommandSpec{
			WorkingDirectory: requestOptionalExistingDirectory("working directory", "inherits scenario"),
			Cmd:              requestCommandLine("command line", true),
		}
	}

	if required {
		_, _ = fmt.Printf("%s:\r\n", description)
		return requestCommand()
	}
	if questionYN(fmt.Sprintf("%s?", description)) {
		return requestCommand()
	}

	return nil
}

func requestEnvVars() map[string]string {
	var envVars map[string]string

	if questionYN("define custom env vars?") {
		envVars = map[string]string{}
		for {
			kv := requestString("K=v", false)
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

func requestScenario() *api.ScenarioSpec {
	return &api.ScenarioSpec{
		Name:             requestString("scenario name", true),
		WorkingDirectory: requestOptionalExistingDirectory("working directory", "inherits benchy's"),
		Env:              requestEnvVars(),
		BeforeAll:        requestCommand("add setup command", false),
		AfterAll:         requestCommand("add teardown command", false),
		BeforeEach:       requestCommand("add before each command", false),
		AfterEach:        requestCommand("add after each command", false),
		Command:          requestCommand("benchmarked command", true),
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
