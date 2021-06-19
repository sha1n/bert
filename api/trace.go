package api

import "time"

// End ends a trace
type End = func(ExecutionInfo, error)

// ID ...
type ID = string

// TraceStream ...
type TraceStream = chan Trace

// Identifiable an abstraction for identifiable objects
type Identifiable interface {
	ID() ID
}

// Tracer a global tracing handler that accumulates trace data and provides access to it.
type Tracer interface {
	Start(i Identifiable) End
	Stream() TraceStream
}

// Trace a single time trace
type Trace interface {
	ID() string
	Elapsed() time.Duration
	System() time.Duration
	User() time.Duration
	Error() error
}

// Stats provides access to statistics. Statistics are not necessarily cached and might be calculated on call.
type Stats interface {
	Min() (time.Duration, error)
	Max() (time.Duration, error)
	Mean() (time.Duration, error)
	Median() (time.Duration, error)
	Percentile(percent float64) (time.Duration, error)
	StdDev() (time.Duration, error)
	ErrorRate() float64
	Count() int
}

// Summary provides access a collection of identifiable statistics.
type Summary interface {
	PerceivedTimeStats(ID) Stats
	SystemTimeStats(ID) Stats
	UserTimeStats(ID) Stats
	IDs() []ID
	Time() time.Time
}
