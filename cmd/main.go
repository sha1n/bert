package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/internal/cli"
	"github.com/sha1n/bert/pkg"
	errorhandling "github.com/sha1n/gommons/pkg/error_handling"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
)

func init() {
	pkg.RegisterInterruptGuard(func(sig os.Signal) {
		termite.NewCursor(os.Stdout).Show()
		doExit(1)
	})

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
			log.Fatal(err)
			exitFn(1)
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
	log.Exit(code)
}

func enableSelfUpdate() bool {
	return DisableSelfUpdate != "true" && runtime.GOOS != "windows"
}
