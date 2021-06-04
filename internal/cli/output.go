package cli

import (
	"errors"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
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

func configureIOContext(cmd *cobra.Command, ctx api.IOContext) api.IOContext {
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

	ctx.DisbaleRichTerminalEffects = !IsExperimentEnabled(cmd, "rich_output")

	log.StandardLogger().SetLevel(level)
	log.StandardLogger().SetOutput(ctx.StderrWriter)

	return ctx
}
