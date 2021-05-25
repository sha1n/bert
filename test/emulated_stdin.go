package test

import (
	"bytes"
	"fmt"
	"strings"
)

type EmulatedStdinReader struct {
	buf           *bytes.Buffer
	lines         []string
	nextLineIndex int
}

func NewEmulatedStdinReader(content string) *EmulatedStdinReader {
	emulatedReader := &EmulatedStdinReader{
		buf:           new(bytes.Buffer),
		lines:         strings.Split(content, "\n"),
		nextLineIndex: 0,
	}

	return emulatedReader
}

func (s *EmulatedStdinReader) Read(buf []byte) (int, error) {
	if s.buf.Len() == 0 && s.nextLineIndex < len(s.lines) {
		s.buf.Reset()
		s.buf.WriteString(fmt.Sprintf("%s\n", s.lines[s.nextLineIndex]))
		s.nextLineIndex++
	}
	return s.buf.Read(buf)
}
