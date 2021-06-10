package report

import (
	"fmt"
	"io"
	"strings"

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
			FormatReportFloatPrecision3(stats.Min),
			FormatReportFloatPrecision3(stats.Max),
			FormatReportFloatPrecision3(stats.Mean),
			FormatReportFloatPrecision3(stats.Median),
			FormatReportFloatPrecision3(func() (float64, error) { return stats.Percentile(90) }),
			FormatReportFloatPrecision3(stats.StdDev),
			FormatReportFloatAsRateInPercents(stats.ErrorRate),
		}); err != nil {
			return err
		}
	}

	return err
}
