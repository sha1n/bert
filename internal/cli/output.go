package cli

import (
	"errors"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	color.NoColor = !termite.Tty
}

var (
	printfRed   = color.New(color.FgRed).Printf
	printRed    = color.New(color.FgRed).Print
	sprintRed   = color.New(color.FgRed).Sprint
	sprintGreen = color.New(color.FgGreen).Sprint
	sprintBold  = color.New(color.Bold).Sprint
)

func configureOutput(cmd *cobra.Command, defaultLogLevel log.Level, ctx api.IOContext) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)
	var level = defaultLogLevel

	if silent && debug {
		CheckUserArgFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		level = log.FatalLevel
	}
	if debug {
		level = log.DebugLevel
	}
	if ctx.Tty {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      true,
		})
	}

	log.StandardLogger().SetLevel(level)
	log.StandardLogger().SetOutput(ctx.StderrWriter)
}
