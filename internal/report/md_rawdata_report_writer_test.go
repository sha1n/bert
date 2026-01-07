package report

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/stretchr/testify/assert"
)

func TestHandleMd(t *testing.T) {
	t1, t2 := twoRandomTraceEvents()

	allRecords := writeMdRawReport(t, true, t1, t2)

	assert.Equal(t, 2+2, len(allRecords))

	// Headers
	assert.Equal(
		t,
		[]string{"Timestamp", "Scenario", "Labels", "Duration", "User Time", "System Time", "Error"},
		allRecords[0],
	)

	assertMdTraceRecord(t, t1, allRecords[2])
	assertMdTraceRecord(t, t2, allRecords[3])
}

func TestHandleMdWithoutHeaders(t *testing.T) {
	t1, t2 := twoRandomTraceEvents()

	allRecords := writeMdRawReport(t, false, t1, t2)

	assert.Equal(t, 2, len(allRecords))

	assertMdTraceRecord(t, t1, allRecords[0])
	assertMdTraceRecord(t, t2, allRecords[1])
}

func writeMdRawReport(t *testing.T, includeHeaders bool, traces ...api.Trace) [][]string {
	expectedRows := len(traces)
	if includeHeaders {
		expectedRows += 2
	}

	readRecords := func(r io.Reader) (records [][]string, err error) {
		records = make([][]string, expectedRows)
		var mdBytes []byte
		if mdBytes, err = io.ReadAll(r); err == nil {
			lines := strings.Split(strings.TrimSpace(string(mdBytes)), "\n")
			for i, line := range lines {
				line = strings.Trim(line, "|")
				cells := strings.Split(line, "|")
				for j, c := range cells {
					cells[j] = strings.TrimSpace(c)
				}
				records[i] = cells
			}
		}
		return records, err
	}

	return writeRawReport(t,
		NewMarkdownStreamReportWriter,
		readRecords,
		includeHeaders,
		traces...,
	)
}

func assertMdTraceRecord(t *testing.T, trace api.Trace, actualRecord []string) {
	expectedLabels := strings.Join(randomLabels, ",")

	_, err := time.Parse(time.RFC3339, actualRecord[0])
	assert.NoError(t, err)
	assert.Equal(t, trace.ID(), actualRecord[1])
	assert.Equal(t, expectedLabels, actualRecord[2])
	assert.Equal(t, FormatReportDuration(func() (time.Duration, error) { return trace.PerceivedTime(), nil }), actualRecord[3])
	assert.Equal(t, FormatReportDuration(func() (time.Duration, error) { return trace.UserCPUTime(), nil }), actualRecord[4])
	assert.Equal(t, FormatReportDuration(func() (time.Duration, error) { return trace.SystemCPUTime(), nil }), actualRecord[5])
	assert.Equal(t, fmt.Sprint(trace.Error() != nil), actualRecord[6])
}
