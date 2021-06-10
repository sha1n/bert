package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/pkg"
	clibtest "github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateMarkdownTableFromWithHeaders(t *testing.T) {
	includeHeaders := true
	lines, summary := generateTestMdReport(t, includeHeaders)
	expectedCellsPerRow := len(SummaryReportHeaders)

	// Verify table structure and dimensions
	assert.Equal(t, 2 /*header + sep*/ +2 /*data*/ +1 /*CRLF*/, len(lines))
	assert.Equal(t, "|Timestamp|Scenario|Samples|Labels|Min|Max|Mean|Median|Percentile 90|StdDev|Errors|", lines[0])
	assert.Equal(t, expectedCellsPerRow, strings.Count(lines[1], "|----"))
	assert.Equal(t, expectedCellsPerRow+1, strings.Count(lines[2], "|"))
	assert.Equal(t, expectedCellsPerRow+1, strings.Count(lines[3], "|"))

	// Verify rows are sorted by ID
	sortedIDs := GetSortedScenarioIds(summary)
	assert.Equal(t, sortedIDs[0], strings.Split(lines[2], "|")[2])
	assert.Equal(t, sortedIDs[1], strings.Split(lines[3], "|")[2])

}

func TestCreateMarkdownTableFromWithNoHeaders(t *testing.T) {
	includeHeaders := false
	lines, summary := generateTestMdReport(t, includeHeaders)
	expectedCellsPerRow := len(SummaryReportHeaders)

	// Verify table structure and dimensions
	assert.Equal(t, 2 /*data*/ +1 /*CRLF*/, len(lines))
	assert.Equal(t, expectedCellsPerRow+1, strings.Count(lines[0], "|"))
	assert.Equal(t, expectedCellsPerRow+1, strings.Count(lines[1], "|"))

	// Verify rows are sorted by ID
	sortedIDs := GetSortedScenarioIds(summary)
	assert.Equal(t, sortedIDs[0], strings.Split(lines[0], "|")[2])
	assert.Equal(t, sortedIDs[1], strings.Split(lines[1], "|")[2])

}

func generateTestMdReport(t *testing.T, includeHeaders bool) ([]string, api.Summary) {
	buf := new(bytes.Buffer)
	writer := buf

	spec := aTwoScenarioSpec()
	t1 := pkg.NewFakeTrace(spec.Scenarios[0].ID(), 1, nil)
	t2 := pkg.NewFakeTrace(spec.Scenarios[1].ID(), 2, errors.New("err2"))

	summary := pkg.NewFakeSummary(t1, t2)
	ctx := api.ReportContext{
		Labels:         clibtest.RandomStrings(),
		IncludeHeaders: includeHeaders,
	}

	reportWriter := NewMarkdownSummaryReportWriter(writer)

	assert.NoError(t, reportWriter(summary, spec, ctx))

	actualMarkdown := buf.String()
	lines := strings.Split(actualMarkdown, "\r\n")

	return lines, summary
}

func aTwoScenarioSpec() api.BenchmarkSpec {
	return api.BenchmarkSpec{
		Executions: 1,
		Scenarios: []api.ScenarioSpec{
			{
				Name: clibtest.RandomString(),
				Command: &api.CommandSpec{
					Cmd: []string{"cmd"},
				},
			},
			{
				Name: clibtest.RandomString(),
				Command: &api.CommandSpec{
					Cmd: []string{"cmd"},
				},
			},
		},
	}
}
