package report

import (
	"io"
	"time"

	"encoding/json"

	"github.com/sha1n/bert/api"
)

// jsonReportWriter a simple human readable test report writer
type jsonReportWriter struct {
	writer io.Writer
}

// NewJSONReportWriter returns a JSON report write handler.
func NewJSONReportWriter(writer io.Writer) api.WriteSummaryReportFn {
	w := jsonReportWriter{
		writer: writer,
	}

	return w.Write
}

func (rw jsonReportWriter) Write(summary api.Summary, config api.BenchmarkSpec, ctx api.ReportContext) (err error) {
	doc := jsonSummaryReportDocument{
		Records: make([]jsonSummaryReportRecord, len(summary.IDs())),
	}

	sortedIds := GetSortedScenarioIds(summary)

	for index, id := range sortedIds {
		stats := summary.PerceivedTimeStats(id)
		userStats := summary.UserTimeStats(id)
		sysStats := summary.SystemTimeStats(id)

		errorRate := float64(stats.ErrorRate())
		doc.Records[index] = jsonSummaryReportRecord{
			Timestamp:  summary.Time().UTC(),
			Name:       id,
			Executions: stats.Count(),
			Labels:     ctx.Labels,
			Min:        floatValueNanos(stats.Min),
			Max:        floatValueNanos(stats.Max),
			Mean:       floatValueNanos(stats.Mean),
			Stddev:     floatValueNanos(stats.StdDev),
			Median:     floatValueNanos(stats.Median),
			P90:        floatValueNanos(func() (time.Duration, error) { return stats.Percentile(90) }),
			User:       floatValueNanos(userStats.Mean),
			System:     floatValueNanos(sysStats.Mean),
			ErrorRate:  &errorRate,
		}
	}

	encoder := json.NewEncoder(rw.writer)
	return encoder.Encode(doc)
}

func floatValueNanos(f func() (time.Duration, error)) (v *int64) {
	value, err := f()
	if err == nil {
		i := value.Nanoseconds()
		v = &i
	}

	return
}

type jsonSummaryReportDocument struct {
	Records []jsonSummaryReportRecord `json:"records,omitempty"`
}

type jsonSummaryReportRecord struct {
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Name       string    `json:"name,omitempty"`
	Executions int       `json:"executions,omitempty"`
	Labels     []string  `json:"labels,omitempty"`
	Min        *int64    `json:"min,omitempty"`
	Max        *int64    `json:"max,omitempty"`
	Mean       *int64    `json:"mean,omitempty"`
	Stddev     *int64    `json:"stddev,omitempty"`
	Median     *int64    `json:"median,omitempty"`
	P90        *int64    `json:"p90,omitempty"`
	User       *int64    `json:"user,omitempty"`
	System     *int64    `json:"system,omitempty"`
	ErrorRate  *float64  `json:"errorRate,omitempty"`
}
