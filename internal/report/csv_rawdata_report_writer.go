package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"encoding/csv"

	log "github.com/sirupsen/logrus"

	"github.com/sha1n/bert/api"
)

// CsvStreamReportWriter a simple human readable test report writer
type CsvStreamReportWriter struct {
	writer *csv.Writer
	ctx    api.ReportContext
}

// NewCsvStreamReportWriter returns a streaming CSV report writer.
func NewCsvStreamReportWriter(writer io.Writer, ctx api.ReportContext) RawDataHandler {
	w := CsvStreamReportWriter{
		writer: csv.NewWriter(writer),
		ctx:    ctx,
	}

	if err := w.writeHeaders(); err != nil {
		log.Error(err)
	}

	return w
}

func (rw CsvStreamReportWriter) writeHeaders() (err error) {
	if rw.ctx.IncludeHeaders {
		err = rw.writer.Write(RawDataReportHeaders)
	}

	return err
}

// Handle handles a real time trace event
func (rw CsvStreamReportWriter) Handle(trace api.Trace) (err error) {
	defer rw.writer.Flush()

	timeStr := FormatDateTime(time.Now(), rw.ctx)
	err = rw.writer.Write([]string{
		timeStr,
		trace.ID(),
		strings.Join(rw.ctx.Labels, ","),
		fmt.Sprintf("%d", trace.Elapsed()),
		fmt.Sprintf("%v", trace.Error() != nil),
	})

	return err
}
