package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sha1n/bert/api"
)

// MarkdownStreamReportWriter a simple human readable test report writer
type MarkdownStreamReportWriter struct {
	writer MarkdownTableWriter
	ctx    api.ReportContext
}

// NewMarkdownStreamReportWriter returns a streaming CSV report writer.
func NewMarkdownStreamReportWriter(writer io.Writer, ctx api.ReportContext) RawDataHandler {
	w := MarkdownStreamReportWriter{
		writer: NewMarkdownTableWriter(writer),
		ctx:    ctx,
	}

	if err := w.writeHeaders(); err != nil {
		log.Error(err)
	}

	return w
}

// Handle handles a real time trace event
func (rw MarkdownStreamReportWriter) Handle(trace api.Trace) (err error) {
	timeStr := FormatDateTime(time.Now(), rw.ctx)
	err = rw.writer.WriteRow([]string{
		timeStr,
		trace.ID(),
		strings.Join(rw.ctx.Labels, ","),
		fmt.Sprintf("%d", trace.PerceivedTime()),
		fmt.Sprintf("%d", trace.UserCPUTime()),
		fmt.Sprintf("%d", trace.SystemCPUTime()),
		fmt.Sprintf("%v", trace.Error() != nil),
	})

	return err
}

func (rw MarkdownStreamReportWriter) writeHeaders() (err error) {
	if rw.ctx.IncludeHeaders {
		err = rw.writer.WriteHeaders(RawDataReportHeaders)
	}

	return err
}
