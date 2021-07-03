package exec

import (
	"errors"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/stretchr/testify/assert"
)

const ZeroErrScenarioID = "noErr"
const SingleErrScenarioID = "err"

func TestErrorRateStat(t *testing.T) {
	summary, _ := generateExampleSummary()

	zeroErrStats := summary.PerceivedTimeStats(ZeroErrScenarioID)
	assert.Equal(t, 0.0, zeroErrStats.ErrorRate())

	singleErrStats := summary.PerceivedTimeStats(SingleErrScenarioID)
	assert.Equal(t, 0.1, singleErrStats.ErrorRate())
}

func TestCount(t *testing.T) {
	summary, expectedCount := generateExampleSummary()

	assertCount := func(stats api.Stats) {
		assert.Equal(t, expectedCount, stats.Count())
	}

	assertCount(summary.PerceivedTimeStats(ZeroErrScenarioID))
	assertCount(summary.UserTimeStats(ZeroErrScenarioID))
	assertCount(summary.SystemTimeStats(ZeroErrScenarioID))
	assertCount(summary.PerceivedTimeStats(SingleErrScenarioID))
	assertCount(summary.UserTimeStats(SingleErrScenarioID))
	assertCount(summary.SystemTimeStats(SingleErrScenarioID))
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
		id:            id,
		perceivedTime: time.Duration(dur),
		error:         err,
	}
}
