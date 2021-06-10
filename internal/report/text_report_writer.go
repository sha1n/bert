package report

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// textReportWriter a simple human readable test report writer
type textReportWriter struct {
	writer  *bufio.Writer
	red     *color.Color
	green   *color.Color
	yellow  *color.Color
	cyan    *color.Color
	magenta *color.Color
	blue    *color.Color
	bold    *color.Color
}

// NewTextReportWriter returns a text report write handler.
func NewTextReportWriter(writer io.Writer, colorsOn bool) api.WriteSummaryReportFn {
	var red, green, yellow, cyan, magenta, blue, bold *color.Color

	if colorsOn {
		red = color.New(color.FgRed)
		green = color.New(color.FgGreen)
		yellow = color.New(color.FgYellow)
		cyan = color.New(color.FgCyan)
		magenta = color.New(color.FgMagenta)
		blue = color.New(color.FgBlue)
		bold = color.New(color.Bold)
	} else {
		red,
			green,
			yellow,
			cyan, magenta, blue,
			bold = color.New(),
			color.New(),
			color.New(),
			color.New(),
			color.New(),
			color.New(),
			color.New()
	}

	w := textReportWriter{
		writer:  bufio.NewWriter(writer),
		red:     red,
		green:   green,
		yellow:  yellow,
		cyan:    cyan,
		magenta: magenta,
		blue:    blue,
		bold:    bold,
	}

	return w.Write
}

func (trw textReportWriter) Write(summary api.Summary, config api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	defer trw.writer.Flush()

	trw.writeNewLine()
	trw.writeTitle(" BENCHMARK SUMMARY")
	trw.writeLabels(ctx.Labels)
	trw.writeDate(summary.Time())
	trw.writeTime(summary.Time())
	trw.writeInt64StatLine("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64StatLine("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writePropertyLine("alternate", config.Alternate)

	trw.writeSeperator()

	sortedIds := GetSortedScenarioIds(summary)

	for _, id := range sortedIds {
		stats := summary.Get(id)

		trw.writeScenarioTitle(id)
		trw.writePropertySec("min", trw.green, stats.Min)
		trw.writePropertySec("mean", trw.cyan, stats.Mean)
		trw.writePropertySec("stddev", trw.blue, stats.StdDev)
		trw.writeNewLine()
		trw.writePropertySec("max", trw.magenta, stats.Max)
		trw.writePropertySec("median", trw.yellow, stats.Median)
		trw.writePropertySec("p90", trw.red, func() (time.Duration, error) { return stats.Percentile(90) })
		trw.writeNewLine()
		trw.writeErrorRateStat("errors", stats.ErrorRate)

		trw.writeSeperator()
	}

	return nil
}

func (trw textReportWriter) writeNewLine() {
	trw.writeString("\n")
}

func (trw textReportWriter) writeSeperator() {
	trw.writeString(fmt.Sprintf("\n%s\n\n", strings.Repeat("-", 63)))
}

func (trw textReportWriter) writeScenarioTitle(name string) {
	trw.writePropertyLine("SCENARIO", trw.yellow.Sprint(name))
}

func (trw textReportWriter) writeTitle(title string) {
	trw.writeString(fmt.Sprintf("%11s\n", title))
}

func (trw textReportWriter) writeLabels(labels []string) {
	trw.writePropertyLine("labels", strings.Join(labels, ","))
}

func (trw textReportWriter) writeDate(time time.Time) {
	trw.writePropertyLine("date", time.Format("Jan 02 2006"))
}

func (trw textReportWriter) writeTime(time time.Time) {
	trw.writePropertyLine("time", time.Format("15:04:05Z07:00"))
}

func (trw textReportWriter) writePropertySec(name string, c *color.Color, f func() (time.Duration, error)) {
	trw.writeProperty(name, FormatReportDuration(f), c)
}

func (trw textReportWriter) writeErrorRateStat(name string, f func() float64) {
	value := f()
	errorRate := int(value * 100)
	if errorRate > 0 {
		trw.writeString(trw.yellow.Sprintf("%11s: %d%%\n", name, errorRate))
	} else {
		trw.writeString(fmt.Sprintf("%11s: %d%%\n", name, errorRate))
	}
}

func (trw textReportWriter) writeInt64StatLine(name string, f func() (int64, error)) {
	trw.writePropertyLine(name, FormatReportInt64(f))
}

func (trw textReportWriter) writePropertyLine(name string, value interface{}) {
	trw.writeProperty(name, value, nil)
	trw.writeNewLine()
}

func (trw textReportWriter) writeProperty(name string, value interface{}, c *color.Color) {
	if c == nil {
		trw.writeString(fmt.Sprintf("%11s: %-7v", name, value))
	} else {
		trw.writeString(c.Sprintf("%11s: ", name))
		trw.writeString(fmt.Sprintf("%-7v", value))
	}
}

func (trw textReportWriter) writeString(str string) {
	_, err := trw.writer.WriteString(str)
	if err != nil {
		log.Error(err)
	}
}
