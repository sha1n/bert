package internal

import (
	"bufio"
	"fmt"
	"strings"
)

type MarkdownTable struct {
	data [][]string
}

type MarkdownTableWriter struct {
	writer *bufio.Writer
}

func NewMarkdownTableWriter(writer *bufio.Writer) *MarkdownTableWriter {
	return &MarkdownTableWriter{
		writer: writer,
	}
}

func (tw *MarkdownTableWriter) WriteHeaders(headers []string) (err error) {
	if err = tw.WriteRow(headers); err == nil {
		_, err = tw.writer.WriteString(fmt.Sprintf("%s|\r\n", strings.Repeat("|----", len(headers))))
	}
	tw.writer.Flush()
	return err
}

func (tw *MarkdownTableWriter) WriteRow(row []string) (err error) {
	// TODO theoretically we need to escape '|' chars
	_, err = tw.writer.WriteString(fmt.Sprintf("|%s|\r\n", strings.Join(row, "|")))
	tw.writer.Flush()
	return err
}

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

func (t *MarkdownTable) SetHeader(col int, header string) *MarkdownTable {
	t.data[0][col] = header
	return t
}

func (t *MarkdownTable) SetData(row, col int, data interface{}) *MarkdownTable {
	return t.setStringData(row, col, fmt.Sprintf("%v", data))
}

func (t *MarkdownTable) SetInt(row, col int, data int) *MarkdownTable {
	return t.setStringData(row, col, fmt.Sprintf("%d", data))
}

func (t *MarkdownTable) SetFloat64(row, col int, data float64) *MarkdownTable {
	return t.setStringData(row, col, fmt.Sprintf("%.3f", data))
}

func (t *MarkdownTable) Write(writer *bufio.Writer) {
	for _, row := range t.data {
		_, _ = writer.WriteString("|")
		for _, col := range row {
			_, _ = writer.WriteString(col)
			_, _ = writer.WriteString("|")
		}
		_, _ = writer.WriteString("\n")
	}

	writer.Flush()
}

func (t *MarkdownTable) setStringData(row, col int, data string) *MarkdownTable {
	row = row + 2
	t.data[row][col] = data
	return t
}
