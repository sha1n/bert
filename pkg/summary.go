package pkg

import (
	"time"

	"github.com/montanaflynn/stats"
	"github.com/sha1n/bert/api"
)

// NewSummary create a new TracerSummary with the specified data.
func NewSummary(tracesByID map[api.ID][]api.Trace) api.Summary {
	summary := &_summary{
		perceivedTimeStats: make(map[api.ID]api.Stats, len(tracesByID)),
		sysCPUTimeStats:    make(map[api.ID]api.Stats, len(tracesByID)),
		userCPUTimeStats:   make(map[api.ID]api.Stats, len(tracesByID)),
		time:               time.Now(),
	}

	for id, traces := range tracesByID {
		perceivedSamples := make([]float64, len(traces))
		systemSamples := make([]float64, len(traces))
		userSamples := make([]float64, len(traces))
		errorCount := 0

		for ti := range traces {
			perceivedSamples[ti] = float64(traces[ti].PerceivedTime().Nanoseconds())
			systemSamples[ti] = float64(traces[ti].SystemCPUTime().Nanoseconds())
			userSamples[ti] = float64(traces[ti].UserCPUTime().Nanoseconds())
			if traces[ti].Error() != nil {
				errorCount++
			}
		}

		summary.perceivedTimeStats[id] = &_stats{
			float64Samples: perceivedSamples,
			errorRate:      float64(errorCount) / float64(len(traces)),
		}
		summary.userCPUTimeStats[id] = &_stats{
			float64Samples: userSamples,
			errorRate:      0,
		}
		summary.sysCPUTimeStats[id] = &_stats{
			float64Samples: systemSamples,
			errorRate:      0,
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
	perceivedTimeStats map[api.ID]api.Stats
	sysCPUTimeStats    map[api.ID]api.Stats
	userCPUTimeStats   map[api.ID]api.Stats
	time               time.Time
}

func (summary *_summary) PerceivedTimeStats(id api.ID) api.Stats {
	return summary.perceivedTimeStats[id]
}

func (summary *_summary) SystemTimeStats(id api.ID) api.Stats {
	return summary.sysCPUTimeStats[id]
}

func (summary *_summary) UserTimeStats(id api.ID) api.Stats {
	return summary.userCPUTimeStats[id]
}

func (summary *_summary) IDs() []api.ID {
	ids := make([]api.ID, 0, len(summary.perceivedTimeStats))
	for k := range summary.perceivedTimeStats {
		ids = append(ids, k)
	}

	return ids
}

func (summary *_summary) Time() time.Time {
	return summary.time
}
