package internal

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aSummaryFor(scenario1, scenario2)

	allRecords := writeCsvReport(t, summary)

	assert.Equal(t, 1+2, len(allRecords))

	// Titles
	assert.Equal(t, []string{"Timestamp", "Scenario", "Min", "Max", "Mean", "Median", "Percentile 90", "StdDev", "Errors"}, allRecords[0])

	expectedTimestamp := summary.Time().Format("2006-01-02T15:04:05Z07:00")

	assertRecord(t, scenario1, summary, expectedTimestamp, allRecords[1])
	assertRecord(t, scenario2, summary, expectedTimestamp, allRecords[2])
}

func assertRecord(t *testing.T, scenario api.Identifiable, summary api.Summary, expectedTimestamp string, actual []string) {
	stats1 := summary.Get(scenario.ID())

	assert.Equal(t, expectedTimestamp, actual[0])
	assert.Equal(t, scenario.ID(), actual[1])
	assert.Equal(t, expectedFloatFormat(stats1.Min), actual[2])
	assert.Equal(t, expectedFloatFormat(stats1.Max), actual[3])
	assert.Equal(t, expectedFloatFormat(stats1.Mean), actual[4])
	assert.Equal(t, expectedFloatFormat(stats1.Median), actual[5])
	assert.Equal(t, expectedFloatFormat(func() (float64, error) { return stats1.Percentile(90) }), actual[6])
	assert.Equal(t, expectedFloatFormat(stats1.StdDev), actual[7])
	assert.Equal(t, expectedRateFormat(stats1.ErrorRate), actual[8])
}

func expectedRateFormat(f func() float64) string {
	errorRate := int(f() * 100)

	return fmt.Sprintf("%d", errorRate)

}

func expectedFloatFormat(f func() (float64, error)) string {
	v, _ := f()
	return fmt.Sprintf("%.3f", v)
}

func writeCsvReport(t *testing.T, summary api.Summary) [][]string {
	buf := new(bytes.Buffer)

	csvWriter := NewCsvReportWriter(bufio.NewWriter(buf))

	csvWriter(summary, nil /* unused */)

	reader := csv.NewReader(buf)
	allRecords, err := reader.ReadAll()

	assert.NoError(t, err)

	return allRecords
}

func aSummaryFor(i1 api.Identifiable, i2 api.Identifiable) api.Summary {
	t := pkg.NewTracer()
	t.Start(i1)(nil)
	t.Start(i2)(nil)
	return t.Summary()
}

type scenario struct {
	id string
}

func (s scenario) ID() string {
	return s.id
}
