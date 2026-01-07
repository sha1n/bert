package report

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// MarkdownTable a markdown table data structure
type MarkdownTable struct {
	data [][]string
}

// MarkdownTableWriter a writer like interface for markdown table data.
type MarkdownTableWriter struct {
	writer *bufio.Writer
}

// NewMarkdownTableWriter creates a new markdown table writer with the specified buffered writer.
func NewMarkdownTableWriter(writer io.Writer) MarkdownTableWriter {
	return MarkdownTableWriter{
		writer: bufio.NewWriter(writer),
	}
}

// WriteHeaders writes table headers line
func (tw MarkdownTableWriter) WriteHeaders(headers []string) (err error) {
	defer func() {
		_ = tw.writer.Flush()
	}()
	if err = tw.WriteRow(headers); err == nil {
		err = tw.writeString(fmt.Sprintf("%s|\r\n", strings.Repeat("|----", len(headers))))
	}

	return err
}

// WriteRow writes a row line
func (tw MarkdownTableWriter) WriteRow(row []string) (err error) {
	// TODO theoretically we need to escape '|' chars
	defer func() {
		_ = tw.writer.Flush()
	}()
	return tw.writeString(fmt.Sprintf("|%s|\r\n", strings.Join(row, "|")))
}

func (tw MarkdownTableWriter) writeString(str string) (err error) {
	_, err = tw.writer.WriteString(str)
	return
}

// NewMarkdownTable creates a new markdown table instance
func NewMarkdownTable(rows, cols int) *MarkdownTable {
	t := new(MarkdownTable)
	rows = rows + 2
	t.data = make([][]string, rows)
	for i := 0; i < rows; i++ {
		t.data[i] = make([]string, cols)
		if i == 1 {
			for j := 0; j < cols; j++ {
				t.data[i][j] = "----"
			}
		}
	}
	return t
}

// SetHeader ...
func (t *MarkdownTable) SetHeader(col int, header string) *MarkdownTable {
	t.data[0][col] = header
	return t
}

// SetData ...
func (t *MarkdownTable) SetData(row, col int, data interface{}) *MarkdownTable {
	return t.setStringData(row, col, fmt.Sprintf("%v", data))
}

// SetInt ...
func (t *MarkdownTable) SetInt(row, col int, data int) *MarkdownTable {
	return t.setStringData(row, col, fmt.Sprintf("%d", data))
}

// SetFloat64 ...
func (t *MarkdownTable) SetFloat64(row, col int, data float64) *MarkdownTable {
	return t.setStringData(
		row, col, fmt.Sprintf("%.3f", data),
	)
}

func (t *MarkdownTable) Write(writer io.Writer) {
	for _, row := range t.data {
		_, _ = io.WriteString(writer, "|")
		for _, col := range row {
			_, _ = io.WriteString(writer, col)
			_, _ = io.WriteString(writer, "|")
		}
		_, _ = io.WriteString(writer, "\n")
	}
}

func (t *MarkdownTable) setStringData(row, col int, data string) *MarkdownTable {
	row = row + 2
	t.data[row][col] = data
	return t
}
