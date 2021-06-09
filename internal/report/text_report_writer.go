package report

import (
	"bufio"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

type format = func(string, ...interface{}) string

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
func NewTextReportWriter(writer *bufio.Writer, colorsOn bool) api.WriteSummaryReportFn {
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
		writer:  writer,
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
	trw.writeNewLine()
	trw.writeTitle(" BENCHMARK SUMMARY")
	trw.writeLabels(ctx.Labels)
	trw.writeDate(summary.Time())
	trw.writeTime(summary.Time())
	trw.writeInt64Stat("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64Stat("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writeBoolProperty("alternate", config.Alternate)

	trw.writeSeperator()

	sortedIds := GetSortedScenarioIds(summary)

	for _, id := range sortedIds {
		stats := summary.Get(id)

		trw.writeScenarioTitle(id)
		trw.writeStatNano2Sec("min", trw.green, stats.Min)
		trw.writeStatNano2Sec("mean", trw.cyan, stats.Mean)
		trw.writeStatNano2Sec("stddev", trw.blue, stats.StdDev)
		trw.writeNewLine()
		trw.writeStatNano2Sec("max", trw.magenta, stats.Max)
		trw.writeStatNano2Sec("median", trw.yellow, stats.Median)
		trw.writeStatNano2Sec("p90", trw.red, func() (float64, error) { return stats.Percentile(90) })
		trw.writeNewLine()
		trw.writeErrorRateStat("errors", stats.ErrorRate)

		trw.writeSeperator()
		trw.writer.Flush()
	}

	return nil
}

func (trw textReportWriter) writeNewLine() {
	trw.writeString("\n")
}

func (trw textReportWriter) writeSeperator() {
	trw.writeString(fmt.Sprintf("\n%s\n\n", strings.Repeat("-", 60)))
}

func (trw textReportWriter) writeScenarioTitle(name string) {
	trw.writeString(fmt.Sprintf("%11s: %s\n", "SCENARIO", trw.yellow.Sprint(name)))
}

func (trw textReportWriter) writeTitle(title string) {
	trw.writeString(fmt.Sprintf("%11s\n", title))
}

func (trw textReportWriter) writeLabels(labels []string) {
	trw.writeString(fmt.Sprintf("%11s: %s\n", "labels", strings.Join(labels, ",")))
}

func (trw textReportWriter) writeDate(time time.Time) {
	timeStr := time.Format("Jan 02 2006")
	trw.writeString(fmt.Sprintf("%11s: %s\n", "date", timeStr))
}

func (trw textReportWriter) writeTime(time time.Time) {
	timeStr := time.Format("15:04:05Z07:00")
	trw.writeString(fmt.Sprintf("%11s: %s\n", "time", timeStr))
}

func (trw textReportWriter) writeStatNano2Sec(name string, c *color.Color, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeString(c.Sprintf("%11s: ", name))
		trw.writeString(fmt.Sprintf("%.3fs", value/math.Pow(10, 9)))
	} else {
		trw.writeStatError(name)
	}
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

func (trw textReportWriter) writeInt64Stat(name string, f func() (int64, error)) {
	value, err := f()
	if err == nil {
		trw.writeString(fmt.Sprintf("%11s: %d\n", name, value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw textReportWriter) writeBoolProperty(name string, value bool) {
	trw.writeString(fmt.Sprintf("%11s: %v\n", name, value))
}

func (trw textReportWriter) writeStatError(name string) {
	trw.writeString(trw.bold.Sprintf("%11s: ", name))
	trw.writeString(trw.red.Sprintf("%s", "ERROR"))
}

func (trw textReportWriter) writeString(str string) {
	_, err := trw.writer.WriteString(str)
	if err != nil {
		log.Error(err)
	}
}
