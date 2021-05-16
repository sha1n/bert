package internal

import (
	"bufio"
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
)

var red = color.New(color.FgRed).Sprintf
var green = color.New(color.FgGreen).Sprintf
var yellow = color.New(color.FgYellow).Sprintf
var bold = color.New(color.Bold).Sprintf

// textReportWriter a simple human readable test report writer
type textReportWriter struct {
	writer *bufio.Writer
}

// NewTextReportWriter returns a text report write handler.
func NewTextReportWriter(writer *bufio.Writer) api.WriteReportFn {
	w := &textReportWriter{
		writer: writer,
	}

	return w.Write
}

func (trw *textReportWriter) Write(ts api.Summary, config *api.BenchmarkSpec) (err error) {
	trw.writeTitle("Benchmark Summary")
	trw.writeInt64Stat("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64Stat("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writeBoolProperty("alternate", config.Alternate)
	trw.writeNewLine()

	for id := range ts.All() {
		stats := ts.Get(id)

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

	return nil
}

func (trw *textReportWriter) writeNewLine() {
	trw.writer.WriteString("\r\n")
}

func (trw *textReportWriter) println(s string) {
	trw.writer.WriteString(fmt.Sprintf("%s\r\n", s))
}

func (trw *textReportWriter) writeTitle(title string) {
	line := strings.Repeat("-", len(title)+2)
	trw.println(line)
	trw.println(green(" %s ", title))
	trw.println(line)
}

func (trw *textReportWriter) writeStatNano2Sec(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%.3fs\r\n", value/math.Pow(10, 9)))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *textReportWriter) writeNumericStat(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%.3f\r\n", value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *textReportWriter) writeErrorRateStat(name string, f func() float64) {
	trw.writeStatTitle(name)

	value := f()
	errorRate := int(value * 100)
	if errorRate > 0 {
		trw.writer.WriteString(yellow("%d%%\r\n", errorRate))
	} else {
		trw.writer.WriteString(fmt.Sprintf("%d%%\r\n", errorRate))
	}
}

func (trw *textReportWriter) writeInt64Stat(name string, f func() (int64, error)) {
	value, err := f()
	if err == nil {
		trw.writeStatTitle(name)
		trw.writer.WriteString(fmt.Sprintf("%d\r\n", value))
	} else {
		trw.writeStatError(name)
	}
}

func (trw *textReportWriter) writeBoolProperty(name string, value bool) {
	trw.writeStatTitle(name)
	trw.writer.WriteString(fmt.Sprintf("%v\r\n", value))
}

func (trw *textReportWriter) writeStatTitle(name string) {
	trw.writer.WriteString(bold("%11s: ", name))
}

func (trw *textReportWriter) writeStatError(name string) {
	trw.writer.WriteString(bold("%11s: ", name))
	trw.writer.WriteString(red("%s\r\n", "ERROR"))
}
