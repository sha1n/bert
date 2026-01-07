package report

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/bert/api"
)

const attentionIndicatorRune = '•'

// textReportWriter a simple human readable test report writer
type textReportWriter struct {
	writer  *bufio.Writer
	red     *color.Color
	green   *color.Color
	yellow  *color.Color
	cyan    *color.Color
	magenta *color.Color
	blue    *color.Color
	hiblue  *color.Color
	bold    *color.Color
}

// NewTextReportWriter returns a text report write handler.
func NewTextReportWriter(writer io.Writer, colorsOn bool) api.WriteSummaryReportFn {
	var red, green, yellow, cyan, magenta, blue, hiblue, bold *color.Color

	if colorsOn {
		red = color.New(color.FgRed)
		green = color.New(color.FgGreen)
		yellow = color.New(color.FgYellow)
		cyan = color.New(color.FgCyan)
		magenta = color.New(color.FgMagenta)
		blue = color.New(color.FgBlue)
		hiblue = color.New(color.FgHiBlue)
		bold = color.New(color.Bold)
	} else {
		red,
			green,
			yellow,
			cyan, magenta, blue, hiblue,
			bold = color.New(),
			color.New(),
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
		hiblue:  hiblue,
		bold:    bold,
	}

	return w.Write
}

func (trw textReportWriter) Write(summary api.Summary, config api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	defer trw.writer.Flush()

	trw.writeNewLine()
	trw.writeTitle(" BENCHMARK SUMMARY")
	trw.writeLabels(ctx.Labels)
	trw.writeDate(summary.Time(), ctx)
	trw.writeTime(summary.Time(), ctx)
	trw.writeInt64StatLine("scenarios", func() (int64, error) { return int64(len(config.Scenarios)), nil })
	trw.writeInt64StatLine("executions", func() (int64, error) { return int64(config.Executions), nil })
	trw.writePropertyLine("alternate", config.Alternate)

	trw.writeSeperator()

	sortedIds := GetSortedScenarioIds(summary)

	for _, id := range sortedIds {
		stats := summary.PerceivedTimeStats(id)
		userStats := summary.UserTimeStats(id)
		sysStats := summary.SystemTimeStats(id)

		trw.writeScenarioTitle(id)
		trw.writeDurationProperty("min", trw.green, stats.Min)
		trw.writeDurationProperty("mean", trw.cyan, stats.Mean)
		trw.writeDurationProperty("median", trw.yellow, stats.Median)
		trw.writeNewLine()

		trw.writeDurationProperty("max", trw.magenta, stats.Max)
		trw.writeDurationProperty("stddev", trw.blue, stats.StdDev)
		trw.writeDurationProperty("p90", trw.red, func() (time.Duration, error) { return stats.Percentile(90) })
		trw.writeNewLine()

		trw.writeDurationProperty("user", trw.hiblue, userStats.Mean)
		trw.writeDurationProperty("system", trw.hiblue, sysStats.Mean)

		trw.writeErrorRateStat("errors", stats.ErrorRate)

		trw.writeNewLine()

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

func (trw textReportWriter) writeDate(time time.Time, ctx api.ReportContext) {
	trw.writePropertyLine("date", FormatDate(time, ctx))
}

func (trw textReportWriter) writeTime(time time.Time, ctx api.ReportContext) {
	trw.writePropertyLine("time", FormatTime(time, ctx))
}

func (trw textReportWriter) writeDurationProperty(name string, c *color.Color, f func() (time.Duration, error)) {
	trw.writeProperty(name, FormatReportDuration(f), c)
}

func (trw textReportWriter) writeErrorRateStat(name string, errorRate func() float64) {
	errorRatePercent := int(errorRate() * 100)
	var attentionIndicator = ""

	if errorRatePercent > 10 {
		attentionIndicator = trw.red.Sprintf("%c", attentionIndicatorRune)
	}

	trw.writeString(fmt.Sprintf("%11s: %d%% %s", name, errorRatePercent, attentionIndicator))
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
		slog.Error(err.Error())
	}
}
