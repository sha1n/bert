package report

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
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

func testTxtSanity(t *testing.T, colorsOn bool) {
	spec := aTwoScenarioSpec()
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
		}{scenario2, time.Second * 2, time.Second * 2, time.Second * 2, true},
	)

	text, lines := writeTxtReport(t, summary, spec, colorsOn)

	assert.Contains(t, text, "date: ")
	assert.Contains(t, text, "time: ")
	assert.Contains(t, lines, fmt.Sprintf("scenarios: %d", len(spec.Scenarios)))
	assert.Contains(t, lines, fmt.Sprintf("executions: %d", spec.Executions))
	assert.Contains(t, lines, fmt.Sprintf("alternate: %v", spec.Alternate))
	assert.Contains(t, lines, fmt.Sprintf("labels: %v", strings.Join(randomLabels, ",")))

	assert.Equal(t, 1, strings.Count(text, "min: 1.0s"))
	assert.Equal(t, 1, strings.Count(text, "min: 2.0s"))

	assert.Equal(t, 1, strings.Count(text, "max: 1.0s"))
	assert.Equal(t, 1, strings.Count(text, "max: 2.0s"))

	assert.Equal(t, 1, strings.Count(text, "median: 1.0s"))
	assert.Equal(t, 1, strings.Count(text, "median: 2.0s"))

	assert.Equal(t, 1, strings.Count(text, "p90: 1.0s"))
	assert.Equal(t, 1, strings.Count(text, "p90: 2.0s"))

	assert.Equal(t, 1, strings.Count(text, "mean: 1.0s"))
	assert.Equal(t, 1, strings.Count(text, "mean: 2.0s"))

	assert.Equal(t, 1, strings.Count(text, "errors: 0%"))
	assert.Equal(t, 1, strings.Count(text, "errors: 100%"))

	if runtime.GOOS != "windows" {
		assert.Equal(t, 1, strings.Count(text, "system: 1"))
		assert.Equal(t, 1, strings.Count(text, "system: 2"))
		assert.Equal(t, 1, strings.Count(text, "user: 1"))
		assert.Equal(t, 1, strings.Count(text, "user: 2"))
	}

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
