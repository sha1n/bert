package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"encoding/csv"

	"github.com/sha1n/benchy/api"
)

// csvReportWriter a simple human readable test report writer
type csvReportWriter struct {
	writer *csv.Writer
}

// NewCsvReportWriter returns a CSV report write handler.
func NewCsvReportWriter(writer io.Writer) api.WriteSummaryReportFn {
	w := csvReportWriter{
		writer: csv.NewWriter(writer),
	}

	return w.Write
}

func (rw csvReportWriter) Write(summary api.Summary, config api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	defer rw.writer.Flush()

	if ctx.IncludeHeaders {
		if err = rw.writer.Write(SummaryReportHeaders); err != nil {
			return err
		}
	}

	timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
	sortedIds := GetSortedScenarioIds(summary)

	for _, id := range sortedIds {
		stats := summary.Get(id)

		if err = rw.writer.Write([]string{
			timeStr,
			id,
			fmt.Sprintf("%d", stats.Count()),
			strings.Join(ctx.Labels, ","),
			FormatReportDurationPlainNanos(stats.Min),
			FormatReportDurationPlainNanos(stats.Max),
			FormatReportDurationPlainNanos(stats.Mean),
			FormatReportDurationPlainNanos(stats.Median),
			FormatReportDurationPlainNanos(func() (time.Duration, error) { return stats.Percentile(90) }),
			FormatReportDurationPlainNanos(stats.StdDev),
			FormatReportFloatAsRateInPercents(stats.ErrorRate),
		}); err != nil {
			return err
		}
	}

	return err
}
