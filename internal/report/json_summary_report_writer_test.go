package report

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_jsonReportWriter_Write(t *testing.T) {
	buffer := new(bytes.Buffer)

	summary := generateReport(t, buffer)
	reportDocument := decodeJSONSummaryReport(t, buffer)

	assert.Equal(t, 2, len(reportDocument.Records))

	for _, record := range reportDocument.Records {
		perceivedStats := summary.PerceivedTimeStats(record.Name)
		userStats := summary.UserTimeStats(record.Name)
		systemStats := summary.SystemTimeStats(record.Name)

		assert.NotNil(t, perceivedStats)
		assert.NotNil(t, userStats)
		assert.NotNil(t, systemStats)

		assert.Equal(t, record.Executions, perceivedStats.Count())
		assert.Equal(t, *record.ErrorRate, perceivedStats.ErrorRate())
		assertStatEqual(t, record.Min, perceivedStats.Min)
		assertStatEqual(t, record.Max, perceivedStats.Max)
		assertStatEqual(t, record.Mean, perceivedStats.Mean)
		assertStatEqual(t, record.Stddev, perceivedStats.StdDev)
		assertStatEqual(t, record.P90, func() (time.Duration, error) { return perceivedStats.Percentile(90) })
		assertStatEqual(t, record.User, userStats.Mean)
		assertStatEqual(t, record.System, systemStats.Mean)
	}
}

func generateReport(t *testing.T, buffer *bytes.Buffer) api.Summary {
	writeFn := NewJSONReportWriter(buffer)
	spec := api.BenchmarkSpec{
		Executions: int(test.RandomUint()),
	}
	summary := aSummary()
	err := writeFn(summary, spec, api.ReportContext{
		Labels:         randomLabels,
		IncludeHeaders: false,
	})

	assert.NoError(t, err)

	return summary
}

func decodeJSONSummaryReport(t *testing.T, buffer *bytes.Buffer) jsonSummaryReportDocument {
	var actual jsonSummaryReportDocument
	err := json.NewDecoder(buffer).Decode(&actual)
	assert.NoError(t, err)

	return actual
}

func assertStatEqual(t *testing.T, actual interface{}, expectedFn func() (time.Duration, error)) {
	expected, err := expectedFn()
	assert.NoError(t, err)

	assert.Equal(t, expected.Nanoseconds(), *(actual.(*int64)))
}

func aSummary() api.Summary {
	return aFakeSummaryFor(
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario{id: "1-id"}, time.Second * time.Duration(test.RandomUint()), time.Second * time.Duration(test.RandomUint()), time.Second * time.Duration(test.RandomUint()), false},
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario{id: "2-id"}, time.Second * time.Duration(test.RandomUint()), time.Second * time.Duration(test.RandomUint()), time.Second * time.Duration(test.RandomUint()), true},
	)

}
