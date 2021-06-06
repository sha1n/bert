package cli

import (
	"errors"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	printfRed   = color.New(color.FgRed).Printf
	printRed    = color.New(color.FgRed).Print
	sprintRed   = color.New(color.FgRed).Sprint
	sprintGreen = color.New(color.FgGreen).Sprint
	sprintBold  = color.New(color.Bold).Sprint
)

func configureOutput(cmd *cobra.Command, ctx IOContext) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)
	var level = log.InfoLevel

	if silent && debug {
		CheckUserArgFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		level = log.PanicLevel
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
