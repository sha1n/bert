package internal

import (
	"bufio"
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/pkg"
)

var red = color.New(color.FgRed).Sprintf
var green = color.New(color.FgGreen).Sprintf
var yellow = color.New(color.FgYellow).Sprintf
var bold = color.New(color.Bold).Sprintf

// ReportWriter an abstraction for an object that benchmark results to any target.
type ReportWriter interface {
	Write(pkg.TracerSummary, *BenchmarkSpec)
}

// TextReportWriter a simple human readable test report writer
type TextReportWriter struct {
	writer *bufio.Writer
}

// NewTextReportWriter creates a new TextReportWriter.
func NewTextReportWriter(writer *bufio.Writer) ReportWriter {
	return &TextReportWriter{
		writer: writer,
	}
}

func (trw *TextReportWriter) Write(ts pkg.TracerSummary, config *BenchmarkSpec) {
	trw.writeTitle("Benchmark Summary")
	trw.writeInt64Stat("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64Stat("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writeBoolProperty("alternate", config.Alternate)
	trw.writeNewLine()

	for id := range ts.AllStats() {
		stats := ts.StatsOf(id)

		title := fmt.Sprintf("Summary of '%s'", id)
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
}

func (trw *TextReportWriter) writeNewLine() {
	trw.writer.WriteString("\r\n")
}

func (trw *TextReportWriter) println(s string) {
	trw.writer.WriteString(fmt.Sprintf("%s\r\n", s))
}

func (trw *TextReportWriter) writeTitle(title string) {
	line := strings.Repeat("-", len(title)+2)
	trw.println(line)
	trw.println(green(" %s ", title))
	trw.println(line)
}

func (trw *TextReportWriter) writeStatNano2Sec(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%.3fs\r\n", value/math.Pow(10, 9)))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *TextReportWriter) writeNumericStat(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%.3f\r\n", value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *TextReportWriter) writeErrorRateStat(name string, f func() float64) {
	trw.writeStatTitle(name)

	value := f()
	errorRate := int(value * 100)
	if errorRate > 0 {
		trw.writer.WriteString(yellow("%d%%\r\n", errorRate))
	} else {
		trw.writer.WriteString(fmt.Sprintf("%d%%\r\n", errorRate))
	}
}

func (trw *TextReportWriter) writeInt64Stat(name string, f func() (int64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%d\r\n", value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *TextReportWriter) writeBoolProperty(name string, value bool) {
	trw.writeStatTitle(name)
	trw.writer.WriteString(fmt.Sprintf("%v\r\n", value))
}

func (trw *TextReportWriter) writeStatTitle(name string) {
	trw.writer.WriteString(bold("%11s: ", name))
}

func (trw *TextReportWriter) writeStatError(name string) {
	trw.writer.WriteString(bold("%11s: ", name))
	trw.writer.WriteString(red("%s\r\n", "ERROR"))
}
