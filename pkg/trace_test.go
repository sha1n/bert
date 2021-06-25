package pkg

import (
	"errors"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_tracer_endFn(t *testing.T) {
	identifiable := api.ScenarioSpec{Name: gommonstest.RandomString()}
	expectedID := identifiable.ID()
	expectedUserTime := time.Duration(gommonstest.RandomUint())
	expectedSysTime := time.Duration(gommonstest.RandomUint())
	expectedPerceivedTime := time.Duration(gommonstest.RandomUint())
	expectedExitCode := int(gommonstest.RandomUint())
	var expectedError error

	tracer := NewTracer(1)

	tracer.Start(identifiable)(
		&api.ExecutionInfo{
			UserTime:      time.Duration(expectedUserTime),
			SystemTime:    time.Duration(expectedSysTime),
			PerceivedTime: time.Duration(expectedPerceivedTime),
			ExitCode:      expectedExitCode,
		},
		expectedError,
	)
	received := <-tracer.Stream()

	assert.Equal(t, expectedPerceivedTime, received.PerceivedTime())
	assert.Equal(t, expectedUserTime, received.UserCPUTime())
	assert.Equal(t, expectedSysTime, received.SystemCPUTime())
	assert.Equal(t, expectedError, received.Error())
	assert.Equal(t, expectedID, received.ID())
}

func Test_tracer_endFn_withError(t *testing.T) {
	identifiable := api.ScenarioSpec{Name: gommonstest.RandomString()}
	expectedID := identifiable.ID()
	expectedError := errors.New(gommonstest.RandomString())
	expectedDuration := time.Nanosecond * 0

	tracer := NewTracer(1)

	tracer.Start(identifiable)(
		nil,
		expectedError,
	)
	received := <-tracer.Stream()

	assert.Equal(t, expectedDuration, received.PerceivedTime())
	assert.Equal(t, expectedDuration, received.UserCPUTime())
	assert.Equal(t, expectedDuration, received.SystemCPUTime())
	assert.Equal(t, expectedError, received.Error())
	assert.Equal(t, expectedID, received.ID())
}
