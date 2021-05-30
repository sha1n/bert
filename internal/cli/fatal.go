package cli

import log "github.com/sirupsen/logrus"

// CheckFatalFn checks the specified error and treats it as fatal if not nil
type CheckFatalFn = func(error)

// CheckBenchmarkInitFatal checks the specified error and treats it as fatal if not nil
func CheckBenchmarkInitFatal(err error) {
	if err != nil {
		log.Errorf("Failed to initialize benchark. Error: %s", err.Error())
		log.Info("Bye!")
		panic(err)
	}
}

// CheckUserArgFatal checks the specified error and treats it as fatal if not nil
func CheckUserArgFatal(err error) {
	if err != nil {
		log.Errorf("Failed to parse program arguments. This is most likely a bug. Error: %s", err.Error())
		log.Info("Bye!")
		panic(err)
	}
}

// CheckFatal checks the specified error and treats it as fatal if not nil
func CheckFatal(err error) {
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		log.Info("Bye!")
		panic(err)
	}
}
