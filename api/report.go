package api

// ReportContext contectual information to propagate to report handlers
type ReportContext struct {
	Labels         []string
	IncludeHeaders bool
}

// WriteReportFn a benchmark report handler
type WriteReportFn = func(Summary, *BenchmarkSpec, *ReportContext) error
