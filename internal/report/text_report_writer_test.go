package report

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

// This test is very loose and is primarily designed to make check that all the critical sections
// and fields of the report exist and the writer don't crash on something basic.
// TODO any clever idea regarding how this can be tighter without too much complexity?
func TestTxtSanity(t *testing.T) {
	spec := aTwoScenarioSpec()
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aSummaryFor(scenario1, scenario2)

	text, lines := writeTxtReport(t, summary, spec)

	assert.NotEqual(t, "", text)
	assert.Equal(
		t,
		31, // header = (3 + 6 + 1) + 2 * (3 + 7 + 1) => 31
		len(lines),
	)

	assert.Contains(t, text, "date: ")
	assert.Contains(t, text, "time: ")
	assert.Contains(t, lines, fmt.Sprintf("scenarios: %d", len(spec.Scenarios)))
	assert.Contains(t, lines, fmt.Sprintf("executions: %d", spec.Executions))
	assert.Contains(t, lines, fmt.Sprintf("alternate: %v", spec.Alternate))
	assert.Contains(t, lines, fmt.Sprintf("labels: %v", strings.Join(randomLabels, ",")))

	assert.Equal(t, 2, strings.Count(text, "min: "))
	assert.Equal(t, 2, strings.Count(text, "max: "))
	assert.Equal(t, 2, strings.Count(text, "median: "))
	assert.Equal(t, 2, strings.Count(text, "p90: "))
	assert.Equal(t, 2, strings.Count(text, "mean: "))
	assert.Equal(t, 2, strings.Count(text, "errors: "))
	assert.Equal(t, 2, strings.Count(text, "stddev: "))

	expectedScenario1Title := expectedTitleFor(scenario1)
	expectedScenario2Title := expectedTitleFor(scenario2)
	assert.Contains(t, lines, expectedScenario1Title)
	assert.Contains(t, lines, expectedScenario1Title)

	assert.Greater(t,
		strings.Index(text, expectedScenario2Title),
		strings.Index(text, expectedScenario1Title),
		"scenario 2 is expected appear after 1, because of sort ordering",
	)

}

func expectedTitleFor(id api.Identifiable) string {
	return fmt.Sprintf("SCENARIO: '%s'", id.ID())
}

func writeTxtReport(t *testing.T, summary api.Summary, spec *api.BenchmarkSpec) (string, []string) {
	buf := new(bytes.Buffer)

	txtWriter := NewTextReportWriter(bufio.NewWriter(buf), false)
	assert.NoError(t,
		txtWriter(
			summary,
			spec,
			&api.ReportContext{
				Labels:         randomLabels,
				IncludeHeaders: false,
			}),
	)

	bytes, err := ioutil.ReadAll(buf)

	assert.NoError(t, err)

	text := strings.TrimSpace(string(bytes))
	lines := strings.Split(text, "\r\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return text, lines
}
