package api

// WriteReportFn a benchmark report handler
type WriteReportFn = func(Summary, *BenchmarkSpec) error
