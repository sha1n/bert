package internal

import (
	"bufio"
	"fmt"
	"strings"

	"encoding/csv"

	"github.com/sha1n/benchy/api"
)

// csvReportWriter a simple human readable test report writer
type csvReportWriter struct {
	writer *csv.Writer
}

// NewCsvReportWriter returns a CSV report write handler.
func NewCsvReportWriter(writer *bufio.Writer) api.WriteSummaryReportFn {
	w := &csvReportWriter{
		writer: csv.NewWriter(writer),
	}

	return w.Write
}

func (rw *csvReportWriter) Write(summary api.Summary, config *api.BenchmarkSpec, ctx *api.ReportContext) (err error) {
	if ctx.IncludeHeaders {
		rw.writer.Write(SummaryReportHeaders)
	}

	timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
	sortedIds := GetSortedScenarioIds(summary)

	for i := range sortedIds {
		id := sortedIds[i]
		stats := summary.Get(id)

		rw.writer.Write([]string{
			timeStr,
			id,
			fmt.Sprintf("%d", stats.Count()),
			strings.Join(ctx.Labels, ","),
			FormatFloat3(stats.Min),
			FormatFloat3(stats.Max),
			FormatFloat3(stats.Mean),
			FormatFloat3(stats.Median),
			FormatFloat3(func() (float64, error) { return stats.Percentile(90) }),
			FormatFloat3(stats.StdDev),
			FormatFloatAsRate(stats.ErrorRate),
		})

		rw.writer.Flush()
	}

	return nil
}
