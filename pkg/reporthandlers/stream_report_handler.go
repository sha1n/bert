package reporthandlers

import (
	"github.com/sha1n/bert/api"
	"github.com/sha1n/bert/pkg/exec"
)

type streamReportHandler struct {
	spec        api.BenchmarkSpec
	ctx         api.ReportContext
	subscriber  exec.StreamSubscriber
	handleFn    exec.HandleFn
	unsubscribe exec.Unsubscribe
}

// NewStreamReportHandler create stream report subscriber.
// Stream report handlers are designed to handle events in real time.
func NewStreamReportHandler(spec api.BenchmarkSpec, ctx api.ReportContext, handleFn exec.HandleFn) api.ReportHandler {
	return &streamReportHandler{
		spec:     spec,
		ctx:      ctx,
		handleFn: handleFn,
	}
}

func (h *streamReportHandler) Subscribe(stream api.TraceStream) {
	h.subscriber = *exec.NewStreamSubscriber(stream, h.handleFn)
	h.unsubscribe = h.subscriber.Subscribe()
}

func (h *streamReportHandler) Finalize() error {
	h.unsubscribe()
	return nil
}
