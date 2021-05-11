package bench

import (
	"github.com/montanaflynn/stats"
)

type TraceStats interface {
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Median() (float64, error)
	Percentile(percent float64) (float64, error)
	Traces() []Trace
}

type traceStats struct {
	traces         []Trace
	float64Samples []float64
}

func (ts *traceStats) Traces() []Trace {
	return ts.traces
}

func (ts *traceStats) Min() (float64, error) {
	return stats.Min(ts.float64Samples)
}

func (ts *traceStats) Max() (float64, error) {
	return stats.Max(ts.float64Samples)
}

func (ts *traceStats) Mean() (float64, error) {
	return stats.Mean(ts.float64Samples)
}

func (ts *traceStats) Median() (float64, error) {
	return stats.Median(ts.float64Samples)
}

func (ts *traceStats) Percentile(percent float64) (float64, error) {
	return stats.Percentile(ts.float64Samples, percent)
}

type TraceSummary interface {
	StatsOf(Id) TraceStats
	AllStats() map[Id]TraceStats
}

type traceSummary struct {
	samples map[Id]TraceStats
}

func (traceSummary *traceSummary) StatsOf(id Id) TraceStats {
	return traceSummary.samples[id]
}

func (traceSummary *traceSummary) AllStats() map[Id]TraceStats {
	return traceSummary.samples
}

func NewTraceSummary(traces map[Id][]Trace) TraceSummary {
	summary := &traceSummary{
		samples: make(map[Id]TraceStats),
	}

	for id, traces := range traces {
		float64Samples := []float64{}
		for ti := range traces {
			float64Samples = append(float64Samples, float64(traces[ti].Elapsed().Nanoseconds()))
		}
		summary.samples[id] = &traceStats{
			traces:         traces,
			float64Samples: float64Samples,
		}
	}

	return summary
}
