package pkg

import (
	"github.com/sha1n/bert/api"
)

type streamReportHandler struct {
	spec        api.BenchmarkSpec
	ctx         api.ReportContext
	subscriber  StreamSubscriber
	handleFn    HandleFn
	unsubscribe Unsubscribe
}

// NewStreamReportHandler create stream report subscriber.
// Stream report handlers are designed to handle events in real time.
func NewStreamReportHandler(spec api.BenchmarkSpec, ctx api.ReportContext, handleFn HandleFn) api.ReportHandler {
	return &streamReportHandler{
		spec:     spec,
		ctx:      ctx,
		handleFn: handleFn,
	}
}

func (h *streamReportHandler) Subscribe(stream api.TraceStream) {
	h.subscriber = *NewStreamSubscriber(stream, h.handleFn)
	h.unsubscribe = h.subscriber.Subscribe()
}

func (h *streamReportHandler) Finalize() error {
	h.unsubscribe()
	return nil
}
