package pkg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
)

const approxSymbol = '≅'

var (
	defaultProgressBarColor = color.New()
	hiYellow                = color.New(color.FgHiYellow)
	yellow                  = color.New(color.FgYellow)
	hiRed                   = color.New(color.FgHiRed)
	red                     = color.New(color.FgRed)
	bold                    = color.New(color.Bold)

	// used to change the color of progress bars when errors are reported
	progressBarErrorColorEscalator = map[*color.Color]*color.Color{
		defaultProgressBarColor: hiYellow,
		hiYellow:                yellow,
		yellow:                  hiRed,
		hiRed:                   red,
		red:                     red,
	}
)

// ProgressView a Listener implementation that uses progress events to render and update
// the terminal visually and in place.
// Combining this implementation with other Stdout writers might break the terminal view.
//
// This implementation is intended to be called from a single thread and is not thread-safe!
// This implementation requires a terminal to be attached to this process. If no terminal is
// attached NewProgressView will panic.
type ProgressView struct {
	matrix           termite.Matrix
	progressInfoByID map[api.ID]*progressInfo
	cursor           termite.Cursor
	cancelHandlers   []context.CancelFunc
	startTime        time.Time
	started          bool
	ended            bool
	eta              etaInfo
}

// NewProgressView creates a new ProgressView for the specified benchmark spec
func NewProgressView(spec api.BenchmarkSpec, termDimensionsFn func() (int, int), ioc api.IOContext) api.Listener {
	scenarioCount := len(spec.Scenarios)
	matrix := termite.NewMatrix(ioc.StdoutWriter, time.Hour)

	rows := matrix.NewRange(
		scenarioCount*4 + // title + progress + empty line per scenario
			2, // ETA + space
	)

	_, termHeight := termDimensionsFn()
	termWidthFn := func() int {
		w, _ := termDimensionsFn()
		return w
	}
	if len(rows)+1 > termHeight {
		yellow.Printf("Your terminal window is too small to dispaly full progress information...\n\n")
		return NewMinimalProgressView(spec, termDimensionsFn, ioc)
	}

	progressInfoByID := make(map[api.ID]*progressInfo, scenarioCount)
	cancelHandlers := make([]context.CancelFunc, scenarioCount)

	nextProgressBarRowIndex := 2

	for i, scenario := range spec.Scenarios {
		formatter := newProgressBarFormatter()
		pBar := termite.NewProgressBar(rows[nextProgressBarRowIndex], spec.Executions, termWidthFn, 59, formatter)
		terminalScaledScenarioName := termite.TruncateString(scenario.Name, termWidthFn()-14)
		rows[nextProgressBarRowIndex-1].Update(fmt.Sprintf("%11s: %s", "SCENARIO", yellow.Sprint(terminalScaledScenarioName)))
		notificationsRowIndex := nextProgressBarRowIndex + 1
		nextProgressBarRowIndex += 4

		tick, cancel, _ := pBar.Start()

		progressInfoByID[scenario.ID()] = &progressInfo{
			minimalProgressInfo: minimalProgressInfo{
				notificationWriter: rows[notificationsRowIndex],
				expectedExecutions: spec.Executions,
			},
			tick:      tick,
			formatter: formatter,
		}
		cancelHandlers[i] = cancel
	}

	return &ProgressView{
		matrix:           matrix,
		eta:              newEtaInfo(rows[len(rows)-1], spec.Alternate, termDimensionsFn),
		progressInfoByID: progressInfoByID,
		cancelHandlers:   cancelHandlers,
		cursor:           termite.NewCursor(ioc.StdoutWriter),
	}
}

// OnBenchmarkStart starts updating view components in the background.
func (l *ProgressView) OnBenchmarkStart() {
	if l.started {
		panic(errors.New("already started"))
	}

	defer l.matrix.UpdateTerminal(true)

	l.started = true
	l.startTime = time.Now()
	restoreLogs := shutOffLogs()
	restoreCursor := l.hideCursor()

	l.cancelHandlers = append(
		l.cancelHandlers,
		restoreCursor,
		restoreLogs,
	)
}

// OnBenchmarkEnd stops all view component updates.
func (l *ProgressView) OnBenchmarkEnd() {
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
func (l *ProgressView) OnScenarioStart(id api.ID) {
	progressInfo := l.progressInfoByID[id]
	progressInfo.lastStartTime = time.Now()
}

// OnScenarioEnd update relevant view components
func (l *ProgressView) OnScenarioEnd(id api.ID) {
	defer l.matrix.UpdateTerminal(true)
	progressInfo := l.progressInfoByID[id]
	elapsed := time.Now().Sub(progressInfo.lastStartTime)
	progressInfo.mean = progressInfo.calculateNewApproxMean(elapsed)
	progressInfo.executions++

	progressInfo.tick(fmt.Sprintf("%-9s", formatDuration(progressInfo.mean)))

	l.eta.update(l.calculateETA(), id)
}

// OnError prints a corresponding error message in the progress info area
func (l *ProgressView) OnError(id api.ID, err error) {
	defer l.matrix.UpdateTerminal(true)
	progressInfo := l.progressInfoByID[id]
	progressInfo.formatter.color = progressBarErrorColorEscalator[l.progressInfoByID[id].formatter.color]
	progressInfo.writeNotification(red.Sprint(err.Error()))
}

// OnMessage prints a corresponding message in the progress info area
func (l *ProgressView) OnMessage(id api.ID, message string) {}

// OnMessagef prints a corresponding message in the progress info area
func (l *ProgressView) OnMessagef(id api.ID, format string, args ...interface{}) {}

func (l *ProgressView) calculateETA() time.Duration {
	var eta time.Duration
	for id := range l.progressInfoByID {
		eta += l.progressInfoByID[id].calculateETA()
	}

	return eta
}

func (l *ProgressView) hideCursor() (restore func()) {
	l.cursor.Hide()
	return func() { l.cursor.Show() }
}

func shutOffLogs() (cancel func()) {
	origLevel := log.GetLevel()
	log.SetLevel(log.FatalLevel)

	return func() { log.SetLevel(origLevel) }
}

type progressBarFormatter struct {
	color *color.Color
}

func newProgressBarFormatter() *progressBarFormatter {
	return &progressBarFormatter{
		color: defaultProgressBarColor,
	}
}

// FormatLeftBorder returns the left border char
func (f *progressBarFormatter) FormatLeftBorder() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatRightBorder returns the right border char
func (f *progressBarFormatter) FormatRightBorder() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatFill returns the fill char
func (f *progressBarFormatter) FormatFill() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatBlank returns the blank char
func (f *progressBarFormatter) FormatBlank() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarBlank)
}

// FormatBlank returns the blank char
func (f *progressBarFormatter) MessageAreaWidth() int {
	return 12
}

type progressInfo struct {
	minimalProgressInfo

	tick      termite.TickMessageFn
	formatter *progressBarFormatter
}

func formatDuration(value time.Duration) string {
	if value.Hours() >= 1 {
		return fmt.Sprintf("%c %.1fh", approxSymbol, value.Hours())
	}
	if value.Minutes() >= 1 {
		return fmt.Sprintf("%c %.1fm", approxSymbol, value.Minutes())
	}
	if value.Seconds() >= 1 {
		return fmt.Sprintf("%c %.1fs", approxSymbol, value.Seconds())
	}
	if value.Milliseconds() >= 1 {
		return fmt.Sprintf("%c %.1fms", approxSymbol, float32(value.Microseconds())/1000)
	}
	if value.Microseconds() >= 1 {
		return fmt.Sprintf("%c %.1fµs", approxSymbol, float32(value.Nanoseconds())/1000)
	}

	return fmt.Sprintf("%c %dns", approxSymbol, value.Nanoseconds())
}
