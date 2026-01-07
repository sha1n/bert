package report

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/sha1n/bert/api"
)

// MarkdownStreamReportWriter a simple human readable test report writer
type MarkdownStreamReportWriter struct {
	writer *bufio.Writer
	ctx    api.ReportContext
}

// NewMarkdownStreamReportWriter returns a new MarkdownStreamReportWriter
func NewMarkdownStreamReportWriter(writer io.Writer, ctx api.ReportContext) RawDataHandler {
	w := &MarkdownStreamReportWriter{
		writer: bufio.NewWriter(writer),
		ctx:    ctx,
	}

	if ctx.IncludeHeaders {
		if err := w.writeHeader(); err != nil {
			slog.Error(err.Error())
		}
	}

	return w
}

// Handle handles a real time trace event
func (rw *MarkdownStreamReportWriter) Handle(trace api.Trace) (err error) {
	_, err = fmt.Fprintf(rw.writer, "| %s | %s | %s | %s | %s | %s | %t |\n",
		FormatDateTime(time.Now(), rw.ctx),
		trace.ID(),
		strings.Join(rw.ctx.Labels, ","),
		FormatReportDuration(func() (time.Duration, error) { return trace.PerceivedTime(), nil }),
		FormatReportDuration(func() (time.Duration, error) { return trace.UserCPUTime(), nil }),
		FormatReportDuration(func() (time.Duration, error) { return trace.SystemCPUTime(), nil }),
		trace.Error() != nil,
	)

	if err == nil {
		err = rw.writer.Flush()
	}

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (rw *MarkdownStreamReportWriter) writeHeader() (err error) {
	_, err = rw.writer.WriteString("| Timestamp | Scenario | Labels | Duration | User Time | System Time | Error |\n")
	if err == nil {
		_, err = rw.writer.WriteString("|-----------|----------|--------|----------|-----------|-------------|-------|\n")
	}

	return err
}
