package api

import "time"

// End ends a trace
type End = func(error)

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
	Error() error
}

// Stats provides access to statistics. Statistics are not necessarily cached and might be calculated on call.
type Stats interface {
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Median() (float64, error)
	Percentile(percent float64) (float64, error)
	StdDev() (float64, error)
	ErrorRate() float64
	Count() int
}

// Summary provides access a cpollection of identifiable statistics.
type Summary interface {
	Get(ID) Stats
	All() map[ID]Stats
	Time() time.Time
}
