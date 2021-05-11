package bench

import (
	"bufio"
	"fmt"
	"math"
	"strings"
)

type ReportWriter interface {
	Write(TraceSummary)
}

type TextReportWriter struct {
	writer *bufio.Writer
}

func NewTextReportWriter(writer *bufio.Writer) ReportWriter {
	return &TextReportWriter{
		writer: writer,
	}
}

func (trw *TextReportWriter) Write(ts TraceSummary) {
	for id := range ts.AllStats() {
		stats := ts.StatsOf(id)

		title := fmt.Sprintf("Summary of '%s'", id)
		trw.writeTitle(title)
		trw.writeIntStat("samples", func() (int64, error) { return int64(len(stats.Traces())), nil })
		trw.writeStatNano2Sec("min (s)", stats.Min)
		trw.writeStatNano2Sec("max (s)", stats.Max)
		trw.writeStatNano2Sec("mean (s)", stats.Mean)
		trw.writeStatNano2Sec("median (s)", stats.Median)
		trw.writeStatNano2Sec("p90 (s)", func() (float64, error) { return stats.Percentile(90) })
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
	line := strings.Repeat("=", len(title)+2)
	trw.println(line)
	trw.println(fmt.Sprintf(" %s ", title))
	trw.println(line)
}

func (trw *TextReportWriter) writeStatNano2Sec(name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		trw.writer.WriteString(fmt.Sprintf("%11s: %.3f\r\n", name, value/math.Pow(10, 9)))
	} else {
		trw.writer.WriteString(fmt.Sprintf("%11s: %s\r\n", name, "ERROR"))
	}
}

func (trw *TextReportWriter) writeIntStat(name string, f func() (int64, error)) {
	value, err := f()
	if err == nil {
		trw.writer.WriteString(fmt.Sprintf("%11s: %d\r\n", name, value))
	} else {
		trw.writer.WriteString(fmt.Sprintf("%11s: %s\r\n", name, "ERROR"))
	}
}
