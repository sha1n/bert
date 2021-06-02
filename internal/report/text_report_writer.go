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
	writer    *bufio.Writer
	fmtRed    func(string, ...interface{}) string
	fmtGreen  func(string, ...interface{}) string
	fmtYellow func(string, ...interface{}) string
	fmtBold   func(string, ...interface{}) string
}

// NewTextReportWriter returns a text report write handler.
func NewTextReportWriter(writer *bufio.Writer, colorsOn bool) api.WriteSummaryReportFn {
	var fmtRed, fmtGreen, fmtYellow, fmtBold format

	if colorsOn {
		fmtRed = color.New(color.FgRed).Sprintf
		fmtGreen = color.New(color.FgGreen).Sprintf
		fmtYellow = color.New(color.FgYellow).Sprintf
		fmtBold = color.New(color.Bold).Sprintf
	} else {
		fmtRed, fmtGreen, fmtYellow, fmtBold = fmt.Sprintf, fmt.Sprintf, fmt.Sprintf, fmt.Sprintf
	}

	w := textReportWriter{
		writer:    writer,
		fmtRed:    fmtRed,
		fmtGreen:  fmtGreen,
		fmtYellow: fmtYellow,
		fmtBold:   fmtBold,
	}

	return w.Write
}

func (trw textReportWriter) Write(summary api.Summary, config api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	trw.writeNewLine()
	trw.writeTitle("BENCHMARK SUMMARY")
	trw.writeLabels(ctx.Labels)
	trw.writeDate(summary.Time())
	trw.writeTime(summary.Time())
	trw.writeInt64Stat("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64Stat("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writeBoolProperty("alternate", config.Alternate)
	trw.writeNewLine()

	sortedIds := GetSortedScenarioIds(summary)

	for _, id := range sortedIds {
		stats := summary.Get(id)

		title := fmt.Sprintf("SCENARIO: '%s'", id)
		trw.writeTitle(title)
		trw.writeStatNano2Sec("min", stats.Min)
		trw.writeStatNano2Sec("max", stats.Max)
		trw.writeStatNano2Sec("mean", stats.Mean)
		trw.writeStatNano2Sec("median", stats.Median)
		trw.writeStatNano2Sec("p90", func() (float64, error) { return stats.Percentile(90) })
		trw.writeStatNano2Sec("stddev", stats.StdDev)
		trw.writeErrorRateStat("errors", stats.ErrorRate)
		trw.writeNewLine()
		trw.writer.Flush()
	}

	return nil
}

func (trw textReportWriter) writeNewLine() {
	trw.writeString("\r\n")
}

func (trw textReportWriter) println(s string) {
	trw.writeString(fmt.Sprintf("%s\r\n", s))
}

func (trw textReportWriter) writeTitle(title string) {
	line := strings.Repeat("-", len(title)+2)
	trw.println(line)
	trw.println(trw.fmtGreen(" %s ", title))
	trw.println(line)
}

func (trw textReportWriter) writeLabels(labels []string) {
	trw.writeStatTitle("labels")
	trw.writeString(fmt.Sprintf("%s\r\n", strings.Join(labels, ",")))
}

func (trw textReportWriter) writeDate(time time.Time) {
	trw.writeStatTitle("date")
	timeStr := time.Format("Jan 02 2006")
	trw.writeString(fmt.Sprintf("%s\r\n", timeStr))
}

func (trw textReportWriter) writeTime(time time.Time) {
	trw.writeStatTitle("time")
	timeStr := time.Format("15:04:05Z07:00")
	trw.writeString(fmt.Sprintf("%s\r\n", timeStr))
}

func (trw textReportWriter) writeStatNano2Sec(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writeString(fmt.Sprintf("%.3fs\r\n", value/math.Pow(10, 9)))
	} else {
		trw.writeStatError(name)
	}
}

func (trw textReportWriter) writeErrorRateStat(name string, f func() float64) {
	trw.writeStatTitle(name)

	value := f()
	errorRate := int(value * 100)
	if errorRate > 0 {
		trw.writeString(trw.fmtYellow("%d%%\r\n", errorRate))
	} else {
		trw.writeString(fmt.Sprintf("%d%%\r\n", errorRate))
	}
}

func (trw textReportWriter) writeInt64Stat(name string, f func() (int64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writeString(fmt.Sprintf("%d\r\n", value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw textReportWriter) writeBoolProperty(name string, value bool) {
	trw.writeStatTitle(name)
	trw.writeString(fmt.Sprintf("%v\r\n", value))
}

func (trw textReportWriter) writeStatTitle(name string) {
	trw.writeString(trw.fmtBold("%11s: ", name))
}

func (trw textReportWriter) writeStatError(name string) {
	trw.writeString(trw.fmtBold("%11s: ", name))
	trw.writeString(trw.fmtRed("%s\r\n", "ERROR"))
}

func (trw textReportWriter) writeString(str string) {
	_, err := trw.writer.WriteString(str)
	if err != nil {
		log.Error(err)
	}
}
