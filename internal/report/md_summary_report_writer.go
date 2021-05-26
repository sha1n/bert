package report

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/sha1n/benchy/api"
)

// mdReportWriter a simple human readable test report writer
type mdReportWriter struct {
	tableWriter *MarkdownTableWriter
}

// NewMarkdownSummaryReportWriter returns a Markdown report write handler.
func NewMarkdownSummaryReportWriter(writer *bufio.Writer) api.WriteSummaryReportFn {
	w := &mdReportWriter{
		tableWriter: NewMarkdownTableWriter(writer),
	}

	return w.Write
}

func (rw *mdReportWriter) Write(summary api.Summary, spec *api.BenchmarkSpec, ctx *api.ReportContext) (err error) {
	if ctx.IncludeHeaders {
		err = rw.tableWriter.WriteHeaders(SummaryReportHeaders)
	}

	if err == nil {
		timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
		sortedIds := GetSortedScenarioIds(summary)
		for _, id := range sortedIds {
			stats := summary.Get(id)

			err = rw.tableWriter.WriteRow([]string{
				timeStr,
				id,
				fmt.Sprint(stats.Count()),
				strings.Join(ctx.Labels, ","),
				FormatReportNanosAsSec3(stats.Min),
				FormatReportNanosAsSec3(stats.Max),
				FormatReportNanosAsSec3(stats.Mean),
				FormatReportNanosAsSec3(stats.Median),
				FormatReportNanosAsSec3(func() (float64, error) { return stats.Percentile(90) }),
				FormatReportNanosAsSec3(stats.StdDev),
				FormatReportFloatAsRate(stats.ErrorRate),
			})
		}

	}

	return err
}
