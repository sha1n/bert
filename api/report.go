package api

// ReportContext contectual information to propagate to report handlers
type ReportContext struct {
	Labels         []string
	IncludeHeaders bool
}

// WriteSummaryReportFn a benchmark report handler
type WriteSummaryReportFn = func(Summary, *BenchmarkSpec, *ReportContext) error

// ReportHandler an async handler
type ReportHandler interface {
	Subscribe(TraceStream)
	Finalize() error
}
