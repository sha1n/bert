package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sha1n/bert/api"
)

// mdReportWriter a simple human readable test report writer
type mdReportWriter struct {
	tableWriter MarkdownTableWriter
}

// NewMarkdownSummaryReportWriter returns a Markdown report write handler.
func NewMarkdownSummaryReportWriter(writer io.Writer) api.WriteSummaryReportFn {
	w := mdReportWriter{
		tableWriter: NewMarkdownTableWriter(writer),
	}

	return w.Write
}

func (rw mdReportWriter) Write(summary api.Summary, spec api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	if ctx.IncludeHeaders {
		err = rw.tableWriter.WriteHeaders(SummaryReportHeaders)
	}

	if err == nil {
		timeStr := FormatDateTime(summary.Time(), ctx)
		sortedIds := GetSortedScenarioIds(summary)
		for _, id := range sortedIds {
			stats := summary.PerceivedTimeStats(id)

			err = rw.tableWriter.WriteRow([]string{
				timeStr,
				id,
				fmt.Sprint(stats.Count()),
				strings.Join(ctx.Labels, ","),
				FormatReportDuration(stats.Min),
				FormatReportDuration(stats.Max),
				FormatReportDuration(stats.Mean),
				FormatReportDuration(stats.Median),
				FormatReportDuration(func() (time.Duration, error) { return stats.Percentile(90) }),
				FormatReportDuration(stats.StdDev),
				FormatReportFloatAsRateInPercents(stats.ErrorRate),
			})
		}

	}

	return err
}
