package cli

import log "github.com/sirupsen/logrus"

func CheckInitFatal(err error) {
	if err != nil {
		log.Errorf("Failed to initialize benchark. Error: %s", err.Error())
		log.Exit(1)
	}
}

func CheckArgFatal(err error) {
	if err != nil {
		log.Errorf("Failed to parse program arguments. This is most likely a bug. Error: %s", err.Error())
		log.Exit(1)
	}
}
