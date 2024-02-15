package cli

import (
	"errors"
	"testing"

	"github.com/sha1n/bert/api"
)

func TestFailFastListener_OnError(t *testing.T) {
	mockID := api.ID("mockID")
	mockError := errors.New("mock error")

	failFastListener := NewFailFastListener(&mockListener{})

	defer func() {
		actual := recover()
		if actual == nil {
			t.Errorf("The code did not panic")
		}

		_, ok := actual.(ExecutionAbortedError)
		if !ok {
			t.Errorf("Unexpected panic type: %T", actual)
		}
	}()

	failFastListener.OnError(mockID, mockError)
}

type mockListener struct{}

// OnBenchmarkEnd implements api.Listener.
func (*mockListener) OnBenchmarkEnd() {
	panic("unimplemented")
}

// OnBenchmarkStart implements api.Listener.
func (*mockListener) OnBenchmarkStart() {
	panic("unimplemented")
}

// OnMessage implements api.Listener.
func (*mockListener) OnMessage(id string, message string) {
	panic("unimplemented")
}

// OnMessagef implements api.Listener.
func (*mockListener) OnMessagef(id string, format string, args ...interface{}) {
	panic("unimplemented")
}

// OnScenarioEnd implements api.Listener.
func (*mockListener) OnScenarioEnd(id string) {
	panic("unimplemented")
}

// OnScenarioStart implements api.Listener.
func (*mockListener) OnScenarioStart(id string) {
	panic("unimplemented")
}

func (m *mockListener) OnError(id api.ID, err error) {}
