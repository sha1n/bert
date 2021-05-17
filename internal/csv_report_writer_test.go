package internal

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"sort"
	"strconv"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	"github.com/stretchr/testify/assert"
)

const ExpectedLinesPerSubject = 7

func TestWrite(t *testing.T) {
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aSummaryFor(scenario1, scenario2)

	allRecords := writeCsvReport(t, summary)

	assert.Equal(t, 1+2*ExpectedLinesPerSubject, len(allRecords))

	// Titles
	assert.Equal(t, []string{"Timestamp", "Scenario", "Stat", "Value"}, allRecords[0])

	expectedTimestamp := summary.Time().Format("2006-01-02T15:04:05Z07:00")
	// scenario 1
	assertRecordRange(t, 1, "1-id", expectedTimestamp, allRecords)
	// scenario 2
	assertRecordRange(t, 8, "2-id", expectedTimestamp, allRecords)
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

func assertRecordRange(t *testing.T, startIndex int, id string, expectedTimestamp string, records [][]string) {
	var statNames []string
	for i := startIndex; i < startIndex+ExpectedLinesPerSubject; i++ {
		assert.Equal(t, 4, len(records[i]))
		assert.Equal(t, expectedTimestamp, records[i][0])
		assert.Equal(t, id, records[i][1])

		_, err := strconv.ParseFloat(records[i][3], 64)
		assert.NoError(t, err)

		statNames = append(statNames, records[i][2])
	}

	expectedStatNames := []string{"errors", "max", "mean", "median", "min", "p90", "stddev"}
	sort.Strings(expectedStatNames)
	sort.Strings(statNames)

	assert.Equal(t, expectedStatNames, statNames)
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
