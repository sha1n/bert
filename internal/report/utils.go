package report

import (
	"fmt"
	"math"
	"sort"

	"github.com/sha1n/benchy/api"
)

// TabularReportDateFormat ...
var TabularReportDateFormat = "2006-01-02T15:04:05Z07:00"

// RawDataReportHeaders ...
var RawDataReportHeaders = []string{"Timestamp", "Scenario", "Labels", "Duration", "Error"}

// SummaryReportHeaders ...
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

// FormatReportFloat3 formats floats for report rendering with 3 digit precision
func FormatReportFloat3(f func() (float64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%.3f", value)
	}

	return "ERROR"
}

// FormatReportNanosAsSec3 formats nano-seconds float as seconds with 3 digits precision for report rendering
func FormatReportNanosAsSec3(f func() (float64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%.3fs", value/math.Pow(10, 9))
	}

	return "ERROR"
}

// FormatReportFloatAsRate formats a float as rate with percent sign, for report rendering
func FormatReportFloatAsRate(f func() float64) string {
	value := f()
	errorRate := int(value * 100)

	return fmt.Sprintf("%d%%", errorRate)
}
