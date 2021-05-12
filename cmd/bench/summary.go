package bench

import (
	"github.com/montanaflynn/stats"
)

type Stats interface {
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Median() (float64, error)
	Percentile(percent float64) (float64, error)
	ErrorRate() float64
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

func (ss *sstats) Median() (float64, error) {
	return stats.Median(ss.float64Samples)
}

func (ss *sstats) Percentile(percent float64) (float64, error) {
	return stats.Percentile(ss.float64Samples, percent)
}

func (ss *sstats) ErrorRate() float64 {
	return ss.errorRate
}

type TracerSummary interface {
	StatsOf(Id) Stats
	AllStats() map[Id]Stats
}

type tracerSummary struct {
	samples map[Id]Stats
}

func (tracerSummary *tracerSummary) StatsOf(id Id) Stats {
	return tracerSummary.samples[id]
}

func (tracerSummary *tracerSummary) AllStats() map[Id]Stats {
	return tracerSummary.samples
}

func NewTracerSummary(tracesById map[Id][]Trace) TracerSummary {
	summary := &tracerSummary{
		samples: make(map[Id]Stats),
	}

	for id, traces := range tracesById {
		float64Samples := []float64{}
		errorCount := 0

		for ti := range traces {
			float64Samples = append(float64Samples, float64(traces[ti].Elapsed().Nanoseconds()))
			if traces[ti].Error() != nil {
				errorCount += 1
			}
		}

		summary.samples[id] = &sstats{
			float64Samples: float64Samples,
			errorRate:      float64(errorCount / len(traces)),
		}
	}

	return summary
}
