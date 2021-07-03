package exec

import (
	"errors"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestNewTraceSink(t *testing.T) {
	stream := make(api.TraceStream)
	sink := NewTraceSink(stream)

	unsubscribe := sink.Subscribe()
	defer unsubscribe()

	trace1, trace2 := radnomTrace(), radnomTrace()

	stream <- trace1
	stream <- trace2

	summary := sink.Summary()

	assert.Equal(t, 2, len(summary.IDs()))
	assert.Contains(t, summary.IDs(), trace1.ID())
	assert.Contains(t, summary.IDs(), trace2.ID())

	assertStatValue(t, trace1.PerceivedTime(), summary.PerceivedTimeStats(trace1.ID()).Min)
	assertStatValue(t, trace1.SystemCPUTime(), summary.SystemTimeStats(trace1.ID()).Mean)
	assertStatValue(t, trace1.UserCPUTime(), summary.UserTimeStats(trace1.ID()).Max)
}

func TestNewTraceSink_Unsubscribe(t *testing.T) {
	stream := make(api.TraceStream, 1)
	sink := NewTraceSink(stream)

	unsubscribe := sink.Subscribe()
	unsubscribe()

	stream <- radnomTrace()

	summary := sink.Summary()

	assert.Equal(t, 0, len(summary.IDs()))
}

func assertStatValue(t *testing.T, expected time.Duration, f func() (time.Duration, error)) {
	value, err := f()

	assert.NoError(t, err)
	assert.Equal(t, expected, value)

}

func radnomTrace() api.Trace {
	var err error
	if test.RandomBool() {
		err = errors.New(test.RandomString())
	}

	return trace{
		id:            test.RandomString(),
		perceivedTime: time.Millisecond * time.Duration(test.RandomUint()),
		usrCPUTime:    time.Millisecond * time.Duration(test.RandomUint()),
		sysCPUTime:    time.Millisecond * time.Duration(test.RandomUint()),
		error:         err,
	}
}
