package internal

import (
	"bufio"
	"fmt"
	"sort"

	"encoding/csv"

	"github.com/sha1n/benchy/api"
)

// textReportWriter a simple human readable test report writer
type csvReportWriter struct {
	writer *csv.Writer
}

// NewCsvReportWriter returns a CSV report write handler.
func NewCsvReportWriter(writer *bufio.Writer) api.WriteReportFn {
	w := &csvReportWriter{
		writer: csv.NewWriter(writer),
	}

	return w.Write
}

func (rw *csvReportWriter) Write(summary api.Summary, config *api.BenchmarkSpec) (err error) {
	rw.writer.Write([]string{"Timestamp", "Scenario", "Stat", "Value"})

	timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
	sortedIds := getSortedScenarioIds(summary)

	for i := range sortedIds {
		id := sortedIds[i]
		stats := summary.Get(id)

		rw.writeStatNano(timeStr, id, "min", stats.Min)
		rw.writeStatNano(timeStr, id, "max", stats.Max)
		rw.writeStatNano(timeStr, id, "mean", stats.Mean)
		rw.writeStatNano(timeStr, id, "median", stats.Median)
		rw.writeStatNano(timeStr, id, "p90", func() (float64, error) { return stats.Percentile(90) })
		rw.writeStatNano(timeStr, id, "stddev", stats.StdDev)
		rw.writeErrorRateStat(timeStr, id, "errors", stats.ErrorRate)
		rw.writer.Flush()
	}

	return nil
}

func getSortedScenarioIds(summary api.Summary) []api.ID {
	statsMap := summary.All()
	sortedIds := make([]api.ID, 0, len(statsMap))
	for k := range statsMap {
		sortedIds = append(sortedIds, k)
	}
	sort.Strings(sortedIds)

	return sortedIds
}

func (rw *csvReportWriter) writeStatNano(timeStr string, id string, name string, f func() (float64, error)) {
	value, err := f()
	if err == nil {
		rw.writer.Write([]string{timeStr, id, name, fmt.Sprintf("%.3f", value)})
	} else {
		rw.writeStatError(id, name)
	}
}

func (rw *csvReportWriter) writeErrorRateStat(timeStr string, id string, name string, f func() float64) {
	value := f()
	errorRate := int(value * 100)

	rw.writer.Write([]string{timeStr, id, name, fmt.Sprintf("%d", errorRate)})
}

func (rw *csvReportWriter) writeStatError(id string, name string) {
	rw.writer.Write([]string{id, name, "ERR"})
}
