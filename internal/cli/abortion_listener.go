package cli

import (
	"fmt"

	"github.com/sha1n/bert/api"
)

// ExecutionAbortedError a marker type for fatal user errors.
// This type of errors is treated differently when user feedback is provided.
type ExecutionAbortedError struct {
	message string
}

func (e ExecutionAbortedError) Error() string {
	return e.message
}

// NewExecutionAbortedError creates a new abortion error with the specified message.
func NewExecutionAbortedError(id api.ID, err error) ExecutionAbortedError {
	return ExecutionAbortedError{
		message: fmt.Sprintf("'%s' reported an error. %s", id, err),
	}
}

type failFastListener struct {
	api.Listener
}

// NewFailFastListener creates a new listener that logs abortion events.
func NewFailFastListener(delegate api.Listener) api.Listener {
	return &failFastListener{Listener: delegate}
}

// OnError logs an error message with the specified ID and error details
func (l failFastListener) OnError(id api.ID, err error) {
	panic(NewExecutionAbortedError(id, err))
}
