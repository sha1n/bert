package report

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
	gommonstest "github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

var randomLabels = append(gommonstest.RandomStrings(), gommonstest.RandomString()) // ensure at least one label

// GetRawDataHandler a provider for a RawDataHandler instance
type GetRawDataHandler = func(io.Writer, api.ReportContext) RawDataHandler

// ParseRecords parses report records from a reader
type ParseRecords = func(io.Reader) ([][]string, error)

func assertRawTraceRecord(t *testing.T, trace api.Trace, actualRecord []string) {
	expectedLabels := strings.Join(randomLabels, ",")

	_, err := time.Parse(time.RFC3339, actualRecord[0])
	assert.NoError(t, err)
	assert.Equal(t, trace.ID(), actualRecord[1])
	assert.Equal(t, expectedLabels, actualRecord[2])
	assert.Equal(t, fmt.Sprint(trace.PerceivedTime().Nanoseconds()), actualRecord[3])
	assert.Equal(t, fmt.Sprint(trace.UserCPUTime().Nanoseconds()), actualRecord[4])
	assert.Equal(t, fmt.Sprint(trace.SystemCPUTime().Nanoseconds()), actualRecord[5])
	assert.Equal(t, fmt.Sprint(trace.Error() != nil), actualRecord[6])
}

func writeRawReport(t *testing.T, getHandler GetRawDataHandler, parseRecords ParseRecords, includeHeaders bool, traces ...api.Trace) [][]string {
	buf := new(bytes.Buffer)
	ctx := api.ReportContext{
		Labels:         randomLabels,
		IncludeHeaders: includeHeaders,
	}

	handler := getHandler(buf, ctx)

	for _, trace := range traces {
		assert.NoError(t, handler.Handle(trace))
	}

	allRecords, err := parseRecords(buf)

	assert.NoError(t, err)

	return allRecords
}
