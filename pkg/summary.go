package pkg

import (
	"time"

	"github.com/montanaflynn/stats"
	"github.com/sha1n/benchy/api"
)

// NewSummary create a new TracerSummary with the specified data.
func NewSummary(tracesByID map[api.ID][]api.Trace) api.Summary {
	summary := &_summary{
		samples: make(map[api.ID]api.Stats),
		time:    time.Now(),
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
			errorRate:      float64(errorCount) / float64(len(traces)),
		}
	}

	return summary
}

type _stats struct {
	float64Samples stats.Float64Data
	errorRate      float64
}

func (s *_stats) Min() (duration time.Duration, err error) {
	return s.nanosStat(stats.Min)
}

func (s *_stats) Max() (duration time.Duration, err error) {
	return s.nanosStat(stats.Max)
}

func (s *_stats) Mean() (duration time.Duration, err error) {
	return s.nanosStat(stats.Mean)
}

func (s *_stats) StdDev() (duration time.Duration, err error) {
	return s.nanosStat(stats.StandardDeviation)
}

func (s *_stats) Median() (duration time.Duration, err error) {
	return s.nanosStat(stats.Median)
}

func (s *_stats) Percentile(percent float64) (duration time.Duration, err error) {
	return s.nanosStat(func(data stats.Float64Data) (float64, error) {
		return stats.Percentile(data, percent)
	})
}

func (s *_stats) ErrorRate() float64 {
	return s.errorRate
}

func (s *_stats) Count() int {
	return len(s.float64Samples)
}

func (s *_stats) nanosStat(f func(stats.Float64Data) (float64, error)) (duration time.Duration, err error) {
	var nanos float64
	if nanos, err = f(s.float64Samples); err == nil {
		duration = time.Duration(nanos) * time.Nanosecond
	}
	return
}

type _summary struct {
	samples map[api.ID]api.Stats
	time    time.Time
}

func (summary *_summary) Get(id api.ID) api.Stats {
	return summary.samples[id]
}

func (summary *_summary) All() map[api.ID]api.Stats {
	return summary.samples
}

func (summary *_summary) Time() time.Time {
	return summary.time
}
