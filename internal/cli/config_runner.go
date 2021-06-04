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
func CreateConfigCommand(ctx api.IOContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Long:  `Interactively walks through a benchmark configuration creation process`,
		Short: `Interactively creates a benchmark config`,
		Run:   createConfigFn(ctx),
	}

	cmd.Flags().StringP(ArgNameOutputFile, "o", "", `output file path. Optional. Writes to stdout by default.`)

	_ = cmd.MarkFlagFilename(ArgNameOutputFile, "yml", "yaml")

	return cmd
}

// createConfigFn returns a function that runs the config tool with the specified IOContext
func createConfigFn(ctx api.IOContext) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx = configureIOContext(cmd, ctx)

		printHints()
		writeCloser := ResolveOutputArg(cmd, ArgNameOutputFile, ctx)
		defer writeCloser.Close()

		spec := api.BenchmarkSpec{
			Executions: int(requestUint("number of executions", true, ctx)),
			Alternate:  requestOptionalBool("alternate executions", false, ctx),
			Scenarios:  requestScenarios(ctx),
		}

		fmt.Print("\r\nWriting your configuration...\r\n\r\n")

		if err := pkg.SaveSpec(spec, writeCloser); err != nil {
			log.Error(err)
			log.Exit(1)
		}
	}
}

func requestScenarios(ctx api.IOContext) []api.ScenarioSpec {
	specs := []api.ScenarioSpec{}

	for {
		specs = append(specs, requestScenario(ctx))
		if !questionYN("add another scenario?", ctx) {
			break
		}
	}

	return specs
}

func requestCommand(description string, required bool, ctx api.IOContext) *api.CommandSpec {
	requestCommand := func() *api.CommandSpec {
		return &api.CommandSpec{
			WorkingDirectory: requestOptionalExistingDirectory("working directory", "inherits scenario", ctx),
			Cmd:              requestCommandLine("command line", true, ctx),
		}
	}

	if required {
		_, _ = fmt.Printf("%s:\r\n", description)
		return requestCommand()
	}
	if questionYN(fmt.Sprintf("%s?", description), ctx) {
		return requestCommand()
	}

	return nil
}

func requestEnvVars(ctx api.IOContext) map[string]string {
	var envVars map[string]string

	if questionYN("define custom env vars?", ctx) {
		envVars = map[string]string{}
		for {
			kv := requestString("K=v", false, ctx)
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

func requestScenario(ctx api.IOContext) api.ScenarioSpec {
	return api.ScenarioSpec{
		Name:             requestString("scenario name", true, ctx),
		WorkingDirectory: requestOptionalExistingDirectory("working directory", "inherits current", ctx),
		Env:              requestEnvVars(ctx),
		BeforeAll:        requestCommand("add setup command", false, ctx),
		AfterAll:         requestCommand("add teardown command", false, ctx),
		BeforeEach:       requestCommand("add before each command", false, ctx),
		AfterEach:        requestCommand("add after each command", false, ctx),
		Command:          requestCommand("benchmarked command", true, ctx),
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
