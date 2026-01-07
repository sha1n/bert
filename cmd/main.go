package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/internal/cli"
	"github.com/sha1n/gommons/pkg/cmd"
	errorhandling "github.com/sha1n/gommons/pkg/error_handling"
)

var (
	// ProgramName : passed from build environment
	ProgramName string
	// Build : passed from build environment
	Build string
	// Version : passed from build environment
	Version string
	// DisableSelfUpdate : passed from build environment to specify that self-update should be disabled
	DisableSelfUpdate string
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
	rootCmd.AddCommand(cmd.CreateShellCompletionScriptGenCommand())
	if enableSelfUpdate() {
		rootCmd.AddCommand(cli.CreateUpdateCommand(Version, ProgramName, ctx))
	}

	if err := rootCmd.Execute(); err != nil {
		doExit(1)
	}
}

func handlePanics(exitFn func(int)) {
	if o := recover(); o != nil {
		if err, ok := o.(cli.FatalUserError); ok {
			slog.Error(err.Error())
			exitFn(1)
		}
		if err, ok := o.(cli.AbortionError); ok {
			slog.Error(err.Error())
			exitFn(0)
		}

		issueURL := errorhandling.GenerateGitHubCreateNewIssueURL(
			"sha1n",
			"bert",
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
	os.Exit(code)
}

func enableSelfUpdate() bool {
	return DisableSelfUpdate != "true" && runtime.GOOS != "windows"
}
