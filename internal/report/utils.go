package report

import (
	"fmt"
	"sort"
	"time"

	"github.com/sha1n/bert/api"
)

// ReportErrorValue ...
const ReportErrorValue = "ERR"

var (
	// RawDataReportHeaders ...
	RawDataReportHeaders = []string{
		"Timestamp", 
		"Scenario", 
		"Labels", 
		"Duration",
		"User Time",
		"System Time",
		"Error",
	}

	// SummaryReportHeaders ...
	SummaryReportHeaders = []string{
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
		"User Time",
		"System Time",
		"Errors",
	}
)

// GetSortedScenarioIds returns a sorted array of scenario IDs for the specified api.Summary
func GetSortedScenarioIds(summary api.Summary) []api.ID {
	sortedIDs := summary.IDs()
	sort.Strings(sortedIDs)

	return sortedIDs
}

// FormatReportDurationPlainNanos formats floats for report rendering with 3 digit precision
func FormatReportDurationPlainNanos(f func() (time.Duration, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprint(value.Nanoseconds())
	}

	return ReportErrorValue
}

// FormatReportInt64 formats int64 values for report rendering
func FormatReportInt64(f func() (int64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%d", value)
	}

	return ReportErrorValue
}

// FormatReportDuration formats nano-seconds float as seconds with 3 digits precision for report rendering
func FormatReportDuration(f func() (time.Duration, error)) string {
	value, err := f()
	if err == nil {
		if value.Hours() >= 1 {
			return fmt.Sprintf("%.1fh", value.Hours())
		}
		if value.Minutes() >= 1 {
			return fmt.Sprintf("%.1fm", value.Minutes())
		}
		if value.Seconds() >= 1 {
			return fmt.Sprintf("%.1fs", value.Seconds())
		}
		if value.Milliseconds() >= 1 {
			return fmt.Sprintf("%.1fms", float32(value.Microseconds())/1000)
		}
		if value.Microseconds() >= 1 {
			return fmt.Sprintf("%.1fµs", float32(value.Nanoseconds())/1000)
		}

		return fmt.Sprintf("%dns", value.Nanoseconds())
	}

	return ReportErrorValue
}

// FormatReportFloatAsRateInPercents formats a float as rate with percent sign, for report rendering
func FormatReportFloatAsRateInPercents(f func() float64) string {
	value := f()
	errorRate := int(value * 100)

	return fmt.Sprintf("%d%%", errorRate)
}

// FormatDateTime formats the specified time using RFC3339 format.
func FormatDateTime(t time.Time, ctx api.ReportContext) string {
	if ctx.UTCDate {
		t = t.UTC()
	}

	return t.Format(time.RFC3339)
}

// FormatDate formats the specified time using RFC3339 format.
func FormatDate(t time.Time, ctx api.ReportContext) string {
	if ctx.UTCDate {
		t = t.UTC()
	}

	return t.Format("Jan 02 2006")
}

// FormatTime formats the specified time using RFC3339 format.
func FormatTime(t time.Time, ctx api.ReportContext) string {
	if ctx.UTCDate {
		t = t.UTC()
	}

	return t.Format("15:04:05")
}
