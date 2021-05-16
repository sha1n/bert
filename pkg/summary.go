package pkg

import (
	"github.com/montanaflynn/stats"
	"github.com/sha1n/benchy/api"
)

// NewSummary create a new TracerSummary with the specified data.
func NewSummary(tracesByID map[api.ID][]api.Trace) api.Summary {
	summary := &_summary{
		samples: make(map[api.ID]api.Stats),
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

		summary.samples[id] = &_stats{
			float64Samples: float64Samples,
			errorRate:      float64(errorCount / len(traces)),
		}
	}

	return summary
}

type _stats struct {
	float64Samples []float64
	errorRate      float64
}

func (s *_stats) Min() (float64, error) {
	return stats.Min(s.float64Samples)
}

func (s *_stats) Max() (float64, error) {
	return stats.Max(s.float64Samples)
}

func (s *_stats) Mean() (float64, error) {
	return stats.Mean(s.float64Samples)
}

func (s *_stats) StdDev() (float64, error) {
	return stats.StandardDeviation(s.float64Samples)
}

func (s *_stats) Median() (float64, error) {
	return stats.Median(s.float64Samples)
}

func (s *_stats) Percentile(percent float64) (float64, error) {
	return stats.Percentile(s.float64Samples, percent)
}

func (s *_stats) ErrorRate() float64 {
	return s.errorRate
}

type _summary struct {
	samples map[api.ID]api.Stats
}

func (summary *_summary) Get(id api.ID) api.Stats {
	return summary.samples[id]
}

func (summary *_summary) All() map[api.ID]api.Stats {
	return summary.samples
}
