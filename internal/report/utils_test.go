package report

import (
	"errors"
	"math/rand"
	"testing"
)

func TestFormatReportFloat3(t *testing.T) {
	tests := []struct {
		name string
		f    func() (float64, error)
		want string
	}{
		{name: "green day", f: func() (float64, error) { return 1, nil }, want: "1.000"},
		{name: "error", f: func() (float64, error) { return 0, errors.New("") }, want: "ERR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatReportFloatPrecision3(tt.f); got != tt.want {
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
		f    func() (float64, error)
		want string
	}{
		{name: "green day below precision threshold", f: func() (float64, error) { return rand.Float64(), nil }, want: "0.000s"},
		{name: "green day above precision threshold", f: func() (float64, error) { return rand.Float64() + 1000000, nil }, want: "0.001s"},
		{name: "error", f: func() (float64, error) { return 0, errors.New("") }, want: "ERR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatReportNanosAsSecPrecision3(tt.f); got != tt.want {
				t.Errorf("FormatReportNanosAsSec3() = %v, want %v", got, tt.want)
			}
		})
	}
}
