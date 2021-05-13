package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/internal"
	"github.com/spf13/cobra"
)

// ProgramName : passed from build environment
var ProgramName string

// Build : passed from build environment
var Build string

// Version : passed from build environment
var Version string

func init() {
	prefixColor := color.New(color.FgWhite, color.Bold)
	log.SetPrefix(prefixColor.Sprint("[BENCHY] "))
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	var rootCmd = &cobra.Command{
		Use: ProgramName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, Version, Build),
		Example:      fmt.Sprintf("%s --config <config file path>", ProgramName),
		SilenceUsage: false,
		Run:          doRun,
	}

	rootCmd.Flags().StringP("config", "c", "", `config file path`)
	rootCmd.Flags().BoolP("pipe-stdout", "", true, `whether or not commands' standard output should be piped to benchy's standard out`)
	rootCmd.Flags().BoolP("pipe-stderr", "", true, `whether or not commands' standard error should be piped to benchy's standard error`)

	cobra.MarkFlagRequired(rootCmd.Flags(), "config")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	_ = rootCmd.Execute()
}

func doRun(cmd *cobra.Command, args []string) {
	error := internal.Run(cmd, args)

	if error != nil {
		log.Println(error.Error())
	}
}
