package report

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

// This test is very loose and is primarily designed to make check that all the critical sections
// and fields of the report exist and the writer don't crash on something basic.
// TODO any clever idea regarding how this can be tighter without too much complexity?
func TestTxtSanityWithoutColors(t *testing.T) {
	testTxtSanity(t, false)
}

func TestTxtSanityWithColors(t *testing.T) {
	testTxtSanity(t, true)
}

func TestWriteStatError(t *testing.T) {
	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	writer := textReportWriter{bufWriter, color.New(), color.New(), color.New(), color.New(), color.New(), color.New(), color.New()}

	name := test.RandomString()
	writer.writeStatError(name)
	bufWriter.Flush()
	assert.Equal(t, name+": ERROR", buf.String())
}

func testTxtSanity(t *testing.T, colorsOn bool) {
	spec := aTwoScenarioSpec()
	var scenario1, scenario2 = scenario{id: "1-id"}, scenario{id: "2-id"}
	summary := aFakeSummaryFor(
		struct {
			id       api.Identifiable
			duration time.Duration
			error    bool
		}{scenario1, time.Second, false},
		struct {
			id       api.Identifiable
			duration time.Duration
			error    bool
		}{scenario2, time.Second * 2, true},
	)

	text, lines := writeTxtReport(t, summary, spec, colorsOn)

	assert.Contains(t, text, "date: ")
	assert.Contains(t, text, "time: ")
	assert.Contains(t, lines, fmt.Sprintf("scenarios: %d", len(spec.Scenarios)))
	assert.Contains(t, lines, fmt.Sprintf("executions: %d", spec.Executions))
	assert.Contains(t, lines, fmt.Sprintf("alternate: %v", spec.Alternate))
	assert.Contains(t, lines, fmt.Sprintf("labels: %v", strings.Join(randomLabels, ",")))

	assert.Equal(t, 1, strings.Count(text, "min: 1.000s"))
	assert.Equal(t, 1, strings.Count(text, "min: 2.000s"))

	assert.Equal(t, 1, strings.Count(text, "max: 1.000s"))
	assert.Equal(t, 1, strings.Count(text, "max: 2.000s"))

	assert.Equal(t, 1, strings.Count(text, "median: 1.000s"))
	assert.Equal(t, 1, strings.Count(text, "median: 2.000s"))

	assert.Equal(t, 1, strings.Count(text, "p90: 1.000s"))
	assert.Equal(t, 1, strings.Count(text, "p90: 2.000s"))

	assert.Equal(t, 1, strings.Count(text, "mean: 1.000s"))
	assert.Equal(t, 1, strings.Count(text, "mean: 2.000s"))

	assert.Equal(t, 1, strings.Count(text, "errors: 0%"))
	assert.Equal(t, 1, strings.Count(text, "errors: 100%"))

	assert.Equal(t, 2, strings.Count(text, "stddev: 0"))

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
	return fmt.Sprintf("SCENARIO: %s", id.ID())
}

func writeTxtReport(t *testing.T, summary api.Summary, spec api.BenchmarkSpec, colorsOn bool) (string, []string) {
	buf := new(bytes.Buffer)

	txtWriter := NewTextReportWriter(buf, colorsOn)
	assert.NoError(t,
		txtWriter(
			summary,
			spec,
			api.ReportContext{
				Labels:         randomLabels,
				IncludeHeaders: false,
			}),
	)

	text := buf.String()
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return text, lines
}
