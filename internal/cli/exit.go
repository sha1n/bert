package cli

import log "github.com/sirupsen/logrus"

type CheckFatalFn = func(error)

func CheckBenchmarkInitFatal(err error) {
	if err != nil {
		log.Errorf("Failed to initialize benchark. Error: %s", err.Error())
		log.Info("Bye!")
		log.Exit(1)
	}
}

func CheckUserArgFatal(err error) {
	if err != nil {
		log.Errorf("Failed to parse program arguments. This is most likely a bug. Error: %s", err.Error())
		log.Info("Bye!")
		log.Exit(1)
	}
}

func CheckFatal(err error) {
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		log.Info("Bye!")
		log.Exit(1)
	}
}
