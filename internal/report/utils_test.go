package report

import (
	"errors"
	"testing"
	"time"

	"github.com/sha1n/bert/api"
)

var now = time.Now()

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

func TestFormatDateTime(t *testing.T) {
	type args struct {
		t   time.Time
		ctx api.ReportContext
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "utc", args: args{t: now, ctx: api.ReportContext{UTCDate: true}}, want: now.UTC().Format(time.RFC3339)},
		{name: "non-utc", args: args{t: now, ctx: api.ReportContext{UTCDate: false}}, want: now.Format(time.RFC3339)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateTime(tt.args.t, tt.args.ctx); got != tt.want {
				t.Errorf("FormatDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	type args struct {
		t   time.Time
		ctx api.ReportContext
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "utc", args: args{t: now, ctx: api.ReportContext{UTCDate: true}}, want: now.UTC().Format("Jan 02 2006")},
		{name: "non-utc", args: args{t: now, ctx: api.ReportContext{UTCDate: false}}, want: now.Format("Jan 02 2006")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.args.t, tt.args.ctx); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	type args struct {
		t   time.Time
		ctx api.ReportContext
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "utc", args: args{t: now, ctx: api.ReportContext{UTCDate: true}}, want: now.UTC().Format("15:04:05")},
		{name: "non-utc", args: args{t: now, ctx: api.ReportContext{UTCDate: false}}, want: now.Format("15:04:05")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.args.t, tt.args.ctx); got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
