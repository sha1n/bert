package report

import (
	"encoding/csv"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	t1, t2 := twoRandomTraceEvents()

	allRecords := writeCsvRawReport(t, true, t1, t2)

	assert.Equal(t, 1+2, len(allRecords))

	// Headers
	assert.Equal(
		t,
		[]string{"Timestamp", "Scenario", "Labels", "Duration", "Error"},
		allRecords[0],
	)

	assertRawTraceRecord(t, t1, allRecords[1])
	assertRawTraceRecord(t, t2, allRecords[2])
}

func TestHandleWithoutHeaders(t *testing.T) {
	t1, t2 := twoRandomTraceEvents()

	allRecords := writeCsvRawReport(t, false, t1, t2)

	assert.Equal(t, 2, len(allRecords))

	assertRawTraceRecord(t, t1, allRecords[0])
	assertRawTraceRecord(t, t2, allRecords[1])
}

func writeCsvRawReport(t *testing.T, includeHeaders bool, traces ...api.Trace) [][]string {
	return writeRawReport(t,
		NewCsvStreamReportWriter,
		func(r io.Reader) ([][]string, error) {
			reader := csv.NewReader(r)
			return reader.ReadAll()
		},
		includeHeaders,
		traces...,
	)
}

func twoRandomTraceEvents() (api.Trace, api.Trace) {
	return pkg.NewFakeTrace(
			gommonstest.RandomString(),
			time.Duration(gommonstest.RandomUint()),
			time.Duration(gommonstest.RandomUint()),
			time.Duration(gommonstest.RandomUint()),
			nil,
		),
		pkg.NewFakeTrace(
			gommonstest.RandomString(),
			time.Duration(gommonstest.RandomUint()),
			time.Duration(gommonstest.RandomUint()),
			time.Duration(gommonstest.RandomUint()),
			errors.New(gommonstest.RandomString()),
		)
}
