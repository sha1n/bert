package pkg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sha1n/benchy/api"
	"github.com/sha1n/termite"
	log "github.com/sirupsen/logrus"
)

var (
	defaultProgressBarColor = color.New()
	hiYellow                = color.New(color.FgHiYellow)
	yellow                  = color.New(color.FgYellow)
	hiRed                   = color.New(color.FgHiRed)
	red                     = color.New(color.FgRed)

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
	matrix               termite.Matrix
	progressHandlersByID map[api.ID]*progressHandlers
	spinner              termite.Spinner
	cursor               termite.Cursor
	cancelHandlers       []context.CancelFunc
	started              bool
	ended                bool
}

// NewProgressView creates a new ProgressView for the specified benchmark spec
func NewProgressView(spec api.BenchmarkSpec, termWidthFn func() int, ioc api.IOContext) api.Listener {
	scenarioCount := len(spec.Scenarios)
	matrix := termite.NewMatrix(ioc.StdoutWriter, time.Millisecond*10)

	rows := matrix.NewRange(
		2 +
			2 + // 2 empty lines
			scenarioCount*3 + // title + progress line  per scenario
			1, // spinner status line
	)
	rows[1].Update("Starting benchmark...")
	rowCount := len(rows)

	progressHandlersByID := make(map[api.ID]*progressHandlers, scenarioCount)
	cancelHandlers := make([]context.CancelFunc, scenarioCount)
	termWidth := termWidthFn()
	progressBarIndex := 4

	for i, scenario := range spec.Scenarios {
		formatter := newProgressBarFormatter()
		pBar := termite.NewProgressBar(rows[progressBarIndex], spec.Executions, termWidth, 40, formatter)
		rows[progressBarIndex-1].Update(yellow.Sprint("- " + scenario.Name))
		progressBarIndex += 3

		tick, cancel, _ := pBar.Start()

		progressHandlersByID[scenario.ID()] = &progressHandlers{
			tick:      tick,
			formatter: formatter,
		}
		cancelHandlers[i] = cancel
	}

	return &ProgressView{
		matrix:               matrix,
		progressHandlersByID: progressHandlersByID,
		cancelHandlers:       cancelHandlers,
		spinner:              termite.NewSpinner(rows[rowCount-1], "Loading...", time.Millisecond*100, termite.DefaultSpinnerFormatter()),
		cursor:               termite.NewCursor(ioc.StdoutWriter),
	}
}

// OnBenchmarkStart starts updating view components in the background.
func (l *ProgressView) OnBenchmarkStart() {
	if l.started {
		panic(errors.New("already started"))
	}

	l.started = true
	restoreLogs := shutOffLogs()
	restoreCursor := l.hideCursor()

	cancelMatrix := l.matrix.Start()
	cancelSpinner, _ := l.spinner.Start()
	l.cancelHandlers = append(
		l.cancelHandlers,
		cancelSpinner,
		cancelMatrix,
		restoreCursor,
		restoreLogs,
		l.cursor.Show,
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

	defer l.cursor.Show()
	defer println()

	_ = l.spinner.Stop("Benchmark finished!")

	// This has to come last, so that the spinner message is updated by the matrix
	for _, cancel := range l.cancelHandlers {
		cancel()
	}
}

// OnScenarioStart does nothing
func (l *ProgressView) OnScenarioStart(id api.ID) {}

// OnScenarioEnd update relevant view components
func (l *ProgressView) OnScenarioEnd(id api.ID) {
	l.progressHandlersByID[id].tick()
}

// OnError prints a corresponding error message in the progress info area
func (l *ProgressView) OnError(id api.ID, err error) {
	l.progressHandlersByID[id].formatter.color = progressBarErrorColorEscalator[l.progressHandlersByID[id].formatter.color]
	_ = l.spinner.SetTitle(color.RedString(err.Error()))
}

// OnMessage prints a corresponding message in the progress info area
func (l *ProgressView) OnMessage(id api.ID, message string) {
	_ = l.spinner.SetTitle(message)
}

// OnMessagef prints a corresponding message in the progress info area
func (l *ProgressView) OnMessagef(id api.ID, format string, args ...interface{}) {
	_ = l.spinner.SetTitle(fmt.Sprintf(format, args...))
}

func (l *ProgressView) hideCursor() (restore func()) {
	l.cursor.Hide()
	cancelCursorHook, _ := registerInterruptGuard(func(os.Signal) {
		l.cursor.Show()
	})

	return func() { cancelCursorHook(); l.cursor.Show() }
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
	return f.color.Sprintf("  %c", termite.DefaultProgressBarFill)
}

// FormatRightBorder returns the right border char
func (f *progressBarFormatter) FormatRightBorder() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarFill)
}

// FormatFill returns the fill char
func (f *progressBarFormatter) FormatFill() string {
	return f.color.Sprintf("%c", termite.DefaultProgressBarFill)
}

type progressHandlers struct {
	tick      termite.TickFn
	formatter *progressBarFormatter
}
