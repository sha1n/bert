package main

import (
	"fmt"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/internal/cli"
	errorhandling "github.com/sha1n/clib/pkg/error_handling"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

var (
	// ProgramName : passed from build environment
	ProgramName string
	// Build : passed from build environment
	Build string
	// Version : passed from build environment
	Version string
)

func main() {
	doRun(doExit)
}

func doRun(exitFn func(int)) {
	defer handlePanics(exitFn)

	ctx := api.NewIOContext()
	rootCmd := cli.NewRootCommand(ProgramName, Version, Build, ctx)

	// Subcommands
	rootCmd.AddCommand(cli.CreateConfigCommand(ctx))
	rootCmd.AddCommand(cli.CreateUpdateCommand(Version, ProgramName, ctx))

	if err := rootCmd.Execute(); err != nil {
		doExit(1)
	}
}

func handlePanics(exitFn func(int)) {
	if o := recover(); o != nil {
		if err, ok := o.(cli.FatalUserError); ok {
			log.Fatal(err)
			exitFn(1)
		}

		issueURL := errorhandling.GenerateGitHubCreateNewIssueURL(
			"sha1n",
			"benchy",
			fmt.Sprintf("Panic Issue (%s, build: %s)", Version, Build),
			fmt.Sprintf("**Error:** %s\n**Stacktrace:**\n```%s```", o, debug.Stack()),
		) + "&labels=bug"

		yellow := color.New(color.FgYellow)
		yellow.Println("\nOh no... Please kindly report this issue by following this URL:")
		fmt.Printf(`

%s

`,
			issueURL,
		)

		exitFn(1)
	}
}

func doExit(code int) {
	log.Exit(code)
}
