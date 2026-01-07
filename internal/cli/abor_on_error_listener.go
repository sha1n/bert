package cli

import (
	"fmt"

	"github.com/sha1n/bert/api"
)

// AbortionError a marker type for fatal user errors.
// This type of errors is treated differently when user feedback is provided.
type AbortionError struct {
	message string
}

func (e AbortionError) Error() string {
	return e.message
}

// NewAbortionError creates a new abortion error with the specified message.
func NewAbortionError(id api.ID, err error) AbortionError {
	return AbortionError{
		message: fmt.Sprintf("'%s' reported an error. %s", id, err),
	}
}

type abortOnErrorListener struct {
	api.Listener
}

// NewAbortOnErrorListener creates a new listener that logs abortion events.
func NewAbortOnErrorListener(delegate api.Listener) api.Listener {
	return &abortOnErrorListener{Listener: delegate}
}

// OnError logs an error message with the specified ID and error details
func (l abortOnErrorListener) OnError(id api.ID, err error) {
	defer panic(NewAbortionError(id, err))
	l.Listener.OnError(id, err)
	l.OnScenarioEnd(id)
}
