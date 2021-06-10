package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/sha1n/benchy/api"
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
		timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
		sortedIds := GetSortedScenarioIds(summary)
		for _, id := range sortedIds {
			stats := summary.Get(id)

			err = rw.tableWriter.WriteRow([]string{
				timeStr,
				id,
				fmt.Sprint(stats.Count()),
				strings.Join(ctx.Labels, ","),
				FormatReportNanosAsSecPrecision3(stats.Min),
				FormatReportNanosAsSecPrecision3(stats.Max),
				FormatReportNanosAsSecPrecision3(stats.Mean),
				FormatReportNanosAsSecPrecision3(stats.Median),
				FormatReportNanosAsSecPrecision3(func() (float64, error) { return stats.Percentile(90) }),
				FormatReportNanosAsSecPrecision3(stats.StdDev),
				FormatReportFloatAsRateInPercents(stats.ErrorRate),
			})
		}

	}

	return err
}
