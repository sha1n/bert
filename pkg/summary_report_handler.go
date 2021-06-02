package pkg

import (
	"github.com/sha1n/benchy/api"
)

type summaryReportHandler struct {
	spec          api.BenchmarkSpec
	ctx           api.ReportContext
	sink          TraceSink
	unsubscribe   Unsubscribe
	writeReportFn api.WriteSummaryReportFn
}

// NewSummaryReportHandler create summary report subscriber
// Summary report handlers typically need to accumulate all the data in order to generate a report.
func NewSummaryReportHandler(spec api.BenchmarkSpec, ctx api.ReportContext, writeReportFn api.WriteSummaryReportFn) api.ReportHandler {
	return &summaryReportHandler{
		spec:          spec,
		ctx:           ctx,
		writeReportFn: writeReportFn,
	}
}

func (h *summaryReportHandler) Subscribe(stream api.TraceStream) {
	h.sink = *NewTraceSink(stream)
	h.unsubscribe = h.sink.Subscribe()
}

func (h *summaryReportHandler) Finalize() error {
	h.unsubscribe()

	return h.writeReportFn(h.sink.Summary(), h.spec, h.ctx)
}
