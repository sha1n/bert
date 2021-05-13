package main

import (
	"fmt"
	"os"

	"github.com/sha1n/benchy/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProgramName : passed from build environment
var ProgramName string

// Build : passed from build environment
var Build string

// Version : passed from build environment
var Version string

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
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
	rootCmd.Flags().BoolP("pipe-stdout", "", true, `redirects external commands standard out to benchy's standard out`)
	rootCmd.Flags().BoolP("pipe-stderr", "", true, `redirects external commands standard error to benchy's standard error`)
	rootCmd.Flags().BoolP("debug", "d", false, `logs extra debug information`)

	cobra.MarkFlagRequired(rootCmd.Flags(), "config")

	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	_ = rootCmd.Execute()
}

func doRun(cmd *cobra.Command, args []string) {
	error := internal.Run(cmd, args)

	if error != nil {
		log.Error(error.Error())
	}
}
