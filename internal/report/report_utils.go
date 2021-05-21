package internal

import (
	"fmt"
	"math"
	"sort"

	"github.com/sha1n/benchy/api"
)

var SummaryReportHeaders = []string{
	"Timestamp",
	"Scenario",
	"Samples",
	"Labels",
	"Min",
	"Max",
	"Mean",
	"Median",
	"Percentile 90",
	"StdDev",
	"Errors",
}

// GetSortedScenarioIds returns a sorted array of scenario IDs for the specified api.Summary
func GetSortedScenarioIds(summary api.Summary) []api.ID {
	statsMap := summary.All()
	sortedIds := make([]api.ID, 0, len(statsMap))
	for k := range statsMap {
		sortedIds = append(sortedIds, k)
	}
	sort.Strings(sortedIds)

	return sortedIds
}

func FormatFloat3(f func() (float64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%.3f", value)
	}

	return "ERROR"
}

func FormatNanosAsSec3(f func() (float64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%.3fs", value/math.Pow(10, 9))
	}

	return "ERROR"
}

func FormatFloatAsRate(f func() float64) string {
	value := f()
	errorRate := int(value * 100)

	return fmt.Sprintf("%d%%", errorRate)
}
