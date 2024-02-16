package cli

import (
	"errors"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/stretchr/testify/mock"
)

func TestAbortOnErrorListener_OnError(t *testing.T) {
	mockID := api.ID("mockID")
	mockError := errors.New("mock error")
	mockListener := new(MockListener)
	mockListener.On("OnError", mockID, mockError).Once()
	mockListener.On("OnScenarioEnd", mockID).Once()
	abortOnErrorListener := NewAbortOnErrorListener(mockListener)

	defer func() {
		actual := recover()
		if actual == nil {
			t.Errorf("The code did not panic")
		}

		_, ok := actual.(AbortionError)
		if !ok {
			t.Errorf("Unexpected panic type: %T", actual)
		}

		mockListener.AssertExpectations(t)
	}()

	abortOnErrorListener.OnError(mockID, mockError)
}

type MockListener struct {
	api.Listener
	mock.Mock
}

func (m *MockListener) OnScenarioEnd(id string) {
	m.Called(id)
}

func (m *MockListener) OnError(id api.ID, err error) {
	m.Called(id, err)
}
