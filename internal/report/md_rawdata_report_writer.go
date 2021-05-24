package report

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sha1n/benchy/api"
)

// MarkdownStreamReportWriter a simple human readable test report writer
type MarkdownStreamReportWriter struct {
	writer *MarkdownTableWriter
	ctx    *api.ReportContext
}

// NewMarkdownStreamReportWriter returns a streaming CSV report writer.
func NewMarkdownStreamReportWriter(writer *bufio.Writer, ctx *api.ReportContext) RawDataHandler {
	w := &MarkdownStreamReportWriter{
		writer: NewMarkdownTableWriter(writer),
		ctx:    ctx,
	}

	if err := w.writeHeaders(); err != nil {
		log.Error(err)
	}

	return w
}

// Handle handles a real time trace event
func (rw *MarkdownStreamReportWriter) Handle(trace api.Trace) (err error) {
	timeStr := time.Now().Format("2006-01-02T15:04:05Z07:00")
	err = rw.writer.WriteRow([]string{
		timeStr,
		trace.ID(),
		strings.Join(rw.ctx.Labels, ","),
		fmt.Sprintf("%d", trace.Elapsed()),
		fmt.Sprintf("%v", trace.Error() != nil),
	})

	return err
}

func (rw *MarkdownStreamReportWriter) writeHeaders() (err error) {
	if rw.ctx.IncludeHeaders {
		err = rw.writer.WriteHeaders(RawDataReportHeaders)
	}

	return err
}
