package pkg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/termite"
)

// MinimalProgressView a Listener implementation that uses minimal screen real-estate to display an ETA
type MinimalProgressView struct {
	matrix           termite.Matrix
	progressInfoByID map[api.ID]*minimalProgressInfo
	started          bool
	ended            bool
	cursor           termite.Cursor
	eta              etaInfo
	alternate        bool
	cancelHandlers   []context.CancelFunc
}

// NewMinimalProgressView creates a new MinimalProgressView for the specified benchmark spec
func NewMinimalProgressView(spec api.BenchmarkSpec, termDimensionsFn func() (int, int), ioc api.IOContext) api.Listener {
	scenarioCount := len(spec.Scenarios)
	progressInfoByID := make(map[api.ID]*minimalProgressInfo, scenarioCount)
	cancelHandlers := []context.CancelFunc{}
	matrix := termite.NewMatrix(ioc.StdoutWriter, time.Hour)
	etaRow := matrix.NewRow()

	for _, scenario := range spec.Scenarios {
		progressInfoByID[scenario.ID()] = &minimalProgressInfo{
			notificationWriter: etaRow,
			expectedExecutions: spec.Executions,
		}
	}

	return &MinimalProgressView{
		matrix:           matrix,
		progressInfoByID: progressInfoByID,
		eta:              newEtaInfo(etaRow, spec.Alternate, termDimensionsFn),
		cursor:           termite.NewCursor(ioc.StdoutWriter),
		alternate:        spec.Alternate,
		cancelHandlers:   cancelHandlers,
	}
}

// OnBenchmarkStart starts updating view components in the background.
func (l *MinimalProgressView) OnBenchmarkStart() {
	if l.started {
		panic(errors.New("already started"))
	}

	l.started = true
	l.cancelHandlers = append(l.cancelHandlers, l.hideCursor(), shutOffLogs())
}

// OnBenchmarkEnd stops all view component updates.
func (l *MinimalProgressView) OnBenchmarkEnd() {
	if !l.started {
		panic(errors.New("not started"))
	}
	if l.ended {
		panic(errors.New("already ended"))
	}
	l.ended = true

	defer l.cursor.Up(1)
	defer l.matrix.UpdateTerminal(false)
	defer l.eta.clear()

	// This has to come last, so that the spinner message is updated by the matrix
	for _, cancel := range l.cancelHandlers {
		cancel()
	}
}

// OnScenarioStart does nothing
func (l *MinimalProgressView) OnScenarioStart(id api.ID) {
	progressInfo := l.progressInfoByID[id]
	progressInfo.lastStartTime = time.Now()
}

// OnScenarioEnd update relevant view components
func (l *MinimalProgressView) OnScenarioEnd(id api.ID) {
	defer l.matrix.UpdateTerminal(true)

	progressInfo := l.progressInfoByID[id]
	elapsed := time.Now().Sub(progressInfo.lastStartTime)
	progressInfo.mean = progressInfo.calculateNewApproxMean(elapsed)
	progressInfo.executions++

	l.eta.update(l.calculateETA(), id)
}

// OnError prints a corresponding error message in the progress info area
func (l *MinimalProgressView) OnError(id api.ID, err error) {
	progressInfo := l.progressInfoByID[id]
	progressInfo.lastError = err
}

// OnMessage prints a corresponding message in the progress info area
func (l *MinimalProgressView) OnMessage(id api.ID, message string) {}

// OnMessagef prints a corresponding message in the progress info area
func (l *MinimalProgressView) OnMessagef(id api.ID, format string, args ...interface{}) {}

func (l *MinimalProgressView) calculateETA() time.Duration {
	var eta time.Duration
	for id := range l.progressInfoByID {
		eta += l.progressInfoByID[id].calculateETA()
	}

	return eta
}

func (l *MinimalProgressView) hideCursor() (restore func()) {
	l.cursor.Hide()
	return func() { l.cursor.Show() }
}

type etaInfo struct {
	writer           io.Writer
	alternate        bool
	termDimensionsFn func() (int, int)
}

func newEtaInfo(writer io.Writer, alternate bool, termDimensionsFn func() (int, int)) (eta etaInfo) {
	eta = etaInfo{
		writer:           writer,
		alternate:        alternate,
		termDimensionsFn: termDimensionsFn,
	}

	defer eta.updateString("pending...")

	return eta
}

func (eta etaInfo) update(dur time.Duration, id api.ID) {
	termWidth, _ := eta.termDimensionsFn()
	if eta.alternate {
		eta.updateStringRaw("%11s: %-9s", "---> ETA", formatDuration(dur))
	} else {
		terminalScaledScenarioName := termite.TruncateString(id, termWidth-36)
		eta.updateStringRaw("%11s: %-9s %s: %s", "---> ETA", formatDuration(dur), "> SCENARIO", yellow.Sprint(terminalScaledScenarioName))
	}
}

func (eta etaInfo) clear() {
	io.WriteString(eta.writer, termite.TermControlEraseLine)
}

func (eta etaInfo) updateString(formattedValue string) {
	eta.updateStringRaw("%11s: %s", "---> ETA", formattedValue)
}

func (eta etaInfo) updateStringRaw(format string, args ...interface{}) {
	io.WriteString(eta.writer, bold.Sprintf(format, args...))
}

type minimalProgressInfo struct {
	notificationWriter io.Writer
	lastStartTime      time.Time
	executions         int
	expectedExecutions int
	mean               time.Duration
	lastError          error
}

func (pi minimalProgressInfo) writeNotification(msg string) {
	io.WriteString(pi.notificationWriter, fmt.Sprintf("%11s  %s", "", msg))
}

func (pi minimalProgressInfo) calculateETA() time.Duration {
	return time.Duration(int64(pi.expectedExecutions-pi.executions) * int64(pi.mean))
}

func (pi minimalProgressInfo) calculateNewApproxMean(elapsed time.Duration) time.Duration {
	if pi.executions == 0 {
		return elapsed
	}
	if pi.executions == pi.expectedExecutions {
		return pi.mean
	}

	meanInNanoseconds := (float64(pi.mean.Nanoseconds())*float64(pi.executions) + float64(elapsed.Nanoseconds())) / float64(pi.executions+1)
	return time.Duration(meanInNanoseconds) * time.Nanosecond

}
