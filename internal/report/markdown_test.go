package report

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownTable(t *testing.T) {
	table := NewMarkdownTable(2, 3)
	table.SetHeader(0, "String")
	table.SetHeader(1, "Int")
	table.SetHeader(2, "Float")
	table.SetData(0, 0, "string")
	table.SetInt(0, 1, 1)
	table.SetFloat64(0, 2, 1.123456789)

	buf := new(bytes.Buffer)

	table.Write(buf)
	mdText := buf.String()

	assert.Equal(t, expectedTableText(), mdText)
}

func expectedTableText() string {
	return `|String|Int|Float|
|----|----|----|
|string|1|1.123|
||||
`
}
