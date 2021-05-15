package pkg

import (
	"github.com/montanaflynn/stats"
)

// Stats provides access to statistics. Statistics are not necessarily cached and might be calculated on call.
type Stats interface {
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Median() (float64, error)
	Percentile(percent float64) (float64, error)
	StdDev() (float64, error)
	ErrorRate() float64
}

// TracerSummary provides access a cpollection of identifiable statistics.
type TracerSummary interface {
	StatsOf(ID) Stats
	AllStats() map[ID]Stats
}

// NewTracerSummary create a new TracerSummary with the specified data.
func NewTracerSummary(tracesByID map[ID][]Trace) TracerSummary {
	summary := &tracerSummary{
		samples: make(map[ID]Stats),
	}

	for id, traces := range tracesByID {
		float64Samples := []float64{}
		errorCount := 0

		for ti := range traces {
			float64Samples = append(float64Samples, float64(traces[ti].Elapsed().Nanoseconds()))
			if traces[ti].Error() != nil {
				errorCount++
			}
		}

		summary.samples[id] = &sstats{
			float64Samples: float64Samples,
			errorRate:      float64(errorCount / len(traces)),
		}
	}

	return summary
}

type sstats struct {
	float64Samples []float64
	errorRate      float64
}

func (ss *sstats) Min() (float64, error) {
	return stats.Min(ss.float64Samples)
}

func (ss *sstats) Max() (float64, error) {
	return stats.Max(ss.float64Samples)
}

func (ss *sstats) Mean() (float64, error) {
	return stats.Mean(ss.float64Samples)
}

func (ss *sstats) StdDev() (float64, error) {
	return stats.StandardDeviation(ss.float64Samples)
}

func (ss *sstats) Median() (float64, error) {
	return stats.Median(ss.float64Samples)
}

func (ss *sstats) Percentile(percent float64) (float64, error) {
	return stats.Percentile(ss.float64Samples, percent)
}

func (ss *sstats) ErrorRate() float64 {
	return ss.errorRate
}

type tracerSummary struct {
	samples map[ID]Stats
}

func (tracerSummary *tracerSummary) StatsOf(id ID) Stats {
	return tracerSummary.samples[id]
}

func (tracerSummary *tracerSummary) AllStats() map[ID]Stats {
	return tracerSummary.samples
}
