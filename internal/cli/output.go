package cli

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func configureOutput(cmd *cobra.Command) {
	silent := GetBool(cmd, ArgNameSilent)
	debug := GetBool(cmd, ArgNameDebug)

	if silent && debug {
		CheckUserArgFatal(errors.New("'--%s' and '--%s' are mutually exclusive"))
	}
	if silent {
		log.StandardLogger().SetLevel(log.PanicLevel)
		log.StandardLogger().SetOutput(os.Stderr)
	}
	if debug {
		log.StandardLogger().SetLevel(log.DebugLevel)
	}
}
