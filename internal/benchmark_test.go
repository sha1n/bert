package internal

import (
	"testing"

	"github.com/sha1n/benchy/pkg"
	"github.com/stretchr/testify/assert"
)

var silentCommandExecutor = NewCommandExecutor(false, false)

func TestRun(t *testing.T) {
	var actualSummary pkg.TracerSummary
	var actualConfig *BenchmarkSpec
	interceptSummary := func(summary pkg.TracerSummary, config *BenchmarkSpec) {
		actualSummary = summary
		actualConfig = config
	}

	err := run("../test/data/benchmark_test_run.yaml", silentCommandExecutor, interceptSummary)

	assert.NoError(t, err)
	assert.NotNil(t, actualSummary)

	assert.Equal(t, 2, len(actualConfig.Scenarios))
	assert.Equal(t, len(actualConfig.Scenarios), len(actualSummary.AllStats()))

	assertFullScenarioStats(t, actualSummary.StatsOf("scenario A"))
	assertFullScenarioStats(t, actualSummary.StatsOf("scenario B"))
}

func TestRunWithMissingConfigFile(t *testing.T) {
	err := run("../test_data/non-existing-file.yaml", silentCommandExecutor, failingWriteSummaryFn(t))

	assert.Error(t, err)
}

func assertFullScenarioStats(t *testing.T, stats pkg.Stats) {
	assert.NotNil(t, stats)
	assert.Equal(t, 0.0, stats.ErrorRate())
	assertStatValue(t, stats.Min)
	assertStatValue(t, stats.Max)
	assertStatValue(t, stats.Mean)
	assertStatValue(t, stats.Median)
}

func assertStatValue(t *testing.T, get func() (float64, error)) {
	value, err := get()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, value, 0.0)
}

func failingWriteSummaryFn(t *testing.T) WriteReport {
	return func(summary pkg.TracerSummary, config *BenchmarkSpec) { t.Fail() }
}
