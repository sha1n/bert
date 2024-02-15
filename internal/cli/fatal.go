package cli

import (
	"fmt"
)

// FatalUserError a marker type for fatal user errors.
// This type of errors is treated differently when user feedback is provided.
type FatalUserError struct {
	message string
}

func (e FatalUserError) Error() string {
	return e.message
}

// NewFatalUserErrorf creates a new fatal user error with the specified message format.
func NewFatalUserErrorf(format string, args ...interface{}) FatalUserError {
	return FatalUserError{
		message: fmt.Sprintf(format, args...),
	}
}

// CheckFatalFn checks the specified error and treats it as fatal if not nil
type CheckFatalFn = func(error)

// CheckBenchmarkInitFatal checks the specified error and treats it as fatal if not nil
func CheckBenchmarkInitFatal(err error) {
	if err != nil {
		panic(NewFatalUserErrorf("Failed to initialize benchmark. Error: %s", err.Error()))
	}
}

// CheckUserArgFatal checks the specified error and treats it as fatal if not nil
func CheckUserArgFatal(err error) {
	if err != nil {
		panic(NewFatalUserErrorf("Failed to parse program arguments. This is most likely a bug. Error: %s", err.Error()))
	}
}

// CheckFatal checks the specified error and treats it as fatal if not nil
func CheckFatal(err error) {
	if err != nil {
		panic(NewFatalUserErrorf("Error: %s", err.Error()))
	}
}
