package main

import (
	"os"

	"github.com/sha1n/benchy/internal/cli"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
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
	defer handlePanics()
	rootCmd := cli.NewRootCommand(ProgramName, Version, Build)

	// Subcommands
	rootCmd.AddCommand(cli.CreateConfigCommand())
	rootCmd.AddCommand(cli.CreateUpdateCommand(Version, ProgramName))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handlePanics() {
	if o := recover(); o != nil {
		if err, ok := o.(error); ok {
			log.Error(err.Error())
		} else {
			log.Error(o)
		}
		log.Exit(1)
	}
}
