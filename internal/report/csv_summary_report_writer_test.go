package report

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aFakeSummaryFor(
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario1, time.Second, time.Second, time.Second, false},
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario2, time.Second, time.Second, time.Second, true},
	)

	allRecords := writeCsvReport(t, summary, true)

	assert.Equal(t, 1+2, len(allRecords))

	// Titles
	assert.Equal(
		t,
		[]string{
			"Timestamp",
			"Scenario",
			"Samples",
			"Labels",
			"Min",
			"Max",
			"Mean",
			"Median",
			"Percentile 90",
			"StdDev",
			"User Time",
			"System Time",
			"Errors",
		},
		allRecords[0],
	)

	expectedTimestamp := summary.Time().Format(time.RFC3339)

	assertRecord(t, scenario1, summary, expectedTimestamp, allRecords[1])
	assertRecord(t, scenario2, summary, expectedTimestamp, allRecords[2])
}

func TestWriteWithNoHeaders(t *testing.T) {
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aFakeSummaryFor(
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario1, time.Second, time.Second, time.Second, false},
		struct {
			id            api.Identifiable
			perceivedTime time.Duration
			userTime      time.Duration
			sysTime       time.Duration
			error         bool
		}{scenario2, time.Second, time.Second, time.Second, true},
	)

	allRecords := writeCsvReport(t, summary, false)

	assert.Equal(t, 2, len(allRecords))

	expectedTimestamp := summary.Time().Format(time.RFC3339)

	assertRecord(t, scenario1, summary, expectedTimestamp, allRecords[0])
	assertRecord(t, scenario2, summary, expectedTimestamp, allRecords[1])
}

func assertRecord(t *testing.T, scenario api.Identifiable, summary api.Summary, expectedTimestamp string, actualRecord []string) {
	stats := summary.PerceivedTimeStats(scenario.ID())
	userStats := summary.UserTimeStats(scenario.ID())
	systemStats := summary.SystemTimeStats(scenario.ID())
	expectedLabels := strings.Join(randomLabels, ",")

	assert.Equal(t, expectedTimestamp, actualRecord[0])
	assert.Equal(t, scenario.ID(), actualRecord[1])
	assert.Equal(t, expectedIntFormat(func() int { return summary.PerceivedTimeStats(scenario.ID()).Count() }), actualRecord[2])
	assert.Equal(t, expectedLabels, actualRecord[3])
	assert.Equal(t, FormatReportDurationPlainNanos(stats.Min), actualRecord[4])
	assert.Equal(t, FormatReportDurationPlainNanos(stats.Max), actualRecord[5])
	assert.Equal(t, FormatReportDurationPlainNanos(stats.Mean), actualRecord[6])
	assert.Equal(t, FormatReportDurationPlainNanos(stats.Median), actualRecord[7])
	assert.Equal(t, FormatReportDurationPlainNanos(func() (time.Duration, error) { return stats.Percentile(90) }), actualRecord[8])
	assert.Equal(t, FormatReportDurationPlainNanos(stats.StdDev), actualRecord[9])
	assert.Equal(t, FormatReportDurationPlainNanos(userStats.Mean), actualRecord[10])
	assert.Equal(t, FormatReportDurationPlainNanos(systemStats.Mean), actualRecord[11])
	assert.Equal(t, expectedRateFormat(stats.ErrorRate), actualRecord[12])
}

func expectedIntFormat(f func() int) string {
	return fmt.Sprintf("%d", f())
}

func expectedRateFormat(f func() float64) string {
	errorRate := int(f() * 100)

	return fmt.Sprintf("%d%%", errorRate)

}

func writeCsvReport(t *testing.T, summary api.Summary, includeHeaders bool) [][]string {
	buf := new(bytes.Buffer)

	csvWriter := NewCsvReportWriter(buf)
	assert.NoError(t,
		csvWriter(
			summary,
			api.BenchmarkSpec{}, /* unused */
			api.ReportContext{
				Labels:         randomLabels,
				IncludeHeaders: includeHeaders,
				UTCDate:        false,
			}),
	)

	reader := csv.NewReader(buf)
	allRecords, err := reader.ReadAll()

	assert.NoError(t, err)

	return allRecords
}

func aFakeSummaryFor(specs ...struct {
	id            api.Identifiable
	perceivedTime time.Duration
	userTime      time.Duration
	sysTime       time.Duration
	error         bool
}) api.Summary {
	traces := []api.Trace{}
	for _, spec := range specs {
		if spec.error {
			traces = append(traces, NewFakeTrace(spec.id.ID(), spec.perceivedTime, spec.userTime, spec.sysTime, errors.New(gommonstest.RandomString())))
		} else {
			traces = append(traces, NewFakeTrace(spec.id.ID(), spec.perceivedTime, spec.userTime, spec.sysTime, nil))
		}
	}

	return NewFakeSummary(traces...)
}

type scenario struct {
	id string
}

func (s scenario) ID() string {
	return s.id
}
