package report

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

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
		[]string{"Timestamp", "Scenario", "Labels", "Duration", "Error"},
		allRecords[0],
	)

	assertRawTraceRecord(t, t1, allRecords[2])
	assertRawTraceRecord(t, t2, allRecords[3])
}

func TestHandleMdWithoutHeaders(t *testing.T) {
	t1, t2 := twoRandomTraceEvents()

	allRecords := writeMdRawReport(t, false, t1, t2)

	assert.Equal(t, 2, len(allRecords))

	assertRawTraceRecord(t, t1, allRecords[0])
	assertRawTraceRecord(t, t2, allRecords[1])
}

func writeMdRawReport(t *testing.T, includeHeaders bool, traces ...api.Trace) [][]string {
	expectedRows := len(traces)
	if includeHeaders {
		expectedRows += 2
	}

	readRecords := func(r io.Reader) (records [][]string, err error) {
		records = make([][]string, expectedRows)
		var mdBytes []byte
		if mdBytes, err = ioutil.ReadAll(r); err == nil {
			lines := strings.Split(strings.TrimSpace(string(mdBytes)), "\r\n")
			for i, line := range lines {
				line = strings.Trim(line, "|")
				records[i] = strings.Split(line, "|")
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
