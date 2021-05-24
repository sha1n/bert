package pkg

import (
	"errors"
	"testing"
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

const ZeroErrScenarioID = "noErr"
const SingleErrScenarioID = "err"

func TestErrorRateStat(t *testing.T) {
	summary, _ := generateExampleSummary()

	zeroErrStats := summary.Get(ZeroErrScenarioID)
	assert.Equal(t, 0.0, zeroErrStats.ErrorRate())

	singleErrStats := summary.Get(SingleErrScenarioID)
	assert.Equal(t, 0.1, singleErrStats.ErrorRate())
}

func TestCount(t *testing.T) {
	summary, size := generateExampleSummary()

	zeroErrStats := summary.Get(ZeroErrScenarioID)
	assert.Equal(t, size, zeroErrStats.Count())

	singleErrStats := summary.Get(SingleErrScenarioID)
	assert.Equal(t, size, singleErrStats.Count())
}

func generateExampleSummary() (api.Summary, int) {
	size := 10
	traces := make(map[api.ID][]api.Trace)

	for i := 1; i <= size; i++ {
		traces[ZeroErrScenarioID] = append(traces[ZeroErrScenarioID], aTraceWith(ZeroErrScenarioID, i, nil))
		if i <= 9 {
			traces[SingleErrScenarioID] = append(traces[SingleErrScenarioID], aTraceWith(SingleErrScenarioID, i, nil))
		}
	}

	traces[SingleErrScenarioID] = append(traces[SingleErrScenarioID], aTraceWith(SingleErrScenarioID, 10, errors.New("test")))

	return NewSummary(traces), size
}

func aTraceWith(id string, dur int, err error) api.Trace {
	return &trace{
		id:      id,
		elapsed: time.Duration(dur),
		error:   err,
	}
}
