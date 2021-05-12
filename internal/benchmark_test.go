package internal

import (
	"testing"

	"github.com/sha1n/benchy/pkg"
	"github.com/stretchr/testify/assert"
)

func Test_run(t *testing.T) {
	var actualSummary pkg.TracerSummary
	var actualConfig *Benchmark
	interceptSummary := func(summary pkg.TracerSummary, config *Benchmark) {
		actualSummary = summary
		actualConfig = config
	}

	err := run("../test_data/benchmark_test_run.yaml", interceptSummary)

	assert.NoError(t, err)
	assert.NotNil(t, actualSummary)

	assert.Equal(t, 2, len(actualConfig.Scenarios))
	assert.Equal(t, len(actualConfig.Scenarios), len(actualSummary.AllStats()))

	assertFullScenarioStats(t, actualSummary.StatsOf("scenario A"))
	assertFullScenarioStats(t, actualSummary.StatsOf("scenario B"))
}

func Test_runWithMissingConfigFile(t *testing.T) {
	err := run("../test_data/non-existing-file.yaml", func(summary pkg.TracerSummary, config *Benchmark) { t.Fail() })

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
