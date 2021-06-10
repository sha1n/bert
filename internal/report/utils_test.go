package report

import (
	"errors"
	"testing"
	"time"
)

func TestFormatReportFloat3(t *testing.T) {
	tests := []struct {
		name string
		f    func() (time.Duration, error)
		want string
	}{
		{name: "green day", f: func() (time.Duration, error) { return 1, nil }, want: "1"},
		{name: "error", f: func() (time.Duration, error) { return 0, errors.New("") }, want: "ERR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatReportDurationPlainNanos(tt.f); got != tt.want {
				t.Errorf("FormatReportFloat3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatReportInt64(t *testing.T) {
	tests := []struct {
		name string
		f    func() (int64, error)
		want string
	}{
		{name: "green day", f: func() (int64, error) { return 1, nil }, want: "1"},
		{name: "error", f: func() (int64, error) { return 0, errors.New("") }, want: "ERR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatReportInt64(tt.f); got != tt.want {
				t.Errorf("FormatReportInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatReportNanosInSecPrecision3(t *testing.T) {
	tests := []struct {
		name string
		f    func() (time.Duration, error)
		want string
	}{
		{name: "green day nanos", f: func() (time.Duration, error) { return time.Duration(1), nil }, want: "1ns"},
		{name: "green day micros", f: func() (time.Duration, error) { return time.Duration(1000), nil }, want: "1.0µs"},
		{name: "green day millis", f: func() (time.Duration, error) { return time.Duration(1000000), nil }, want: "1.0ms"},
		{name: "green day seconds", f: func() (time.Duration, error) { return time.Duration(1000000000), nil }, want: "1.0s"},
		{name: "green day minutes", f: func() (time.Duration, error) { return time.Duration(60 * 1000000000), nil }, want: "1.0m"},
		{name: "green day hours", f: func() (time.Duration, error) { return time.Duration(60 * 60000000000), nil }, want: "1.0h"},
		{name: "error", f: func() (time.Duration, error) { return 0, errors.New("") }, want: "ERR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatReportDuration(tt.f); got != tt.want {
				t.Errorf("FormatReportNanosAsSec3() = %v, want %v", got, tt.want)
			}
		})
	}
}
