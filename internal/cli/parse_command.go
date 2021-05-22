package cli

import (
	"unicode"
)

const singleQuote = '\''
const doubleQuote = '"'

func parseCommand(cmd string) []string {
	command := []string{}
	segment := ""
	var usedQuote rune
	escapeNext := false
	quote := false
	for _, c := range cmd {
		if escapeNext {
			escapeNext = false
			segment += string(c)
			continue
		}

		if c == singleQuote || c == doubleQuote {
			if quote && c != usedQuote {
				segment += string(c)

			} else if len(segment) == 0 {
				// mark a new quoted segment
				usedQuote = c
			}

			if c == usedQuote {
				quote = !quote
			}

			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if unicode.IsSpace(c) && !quote {
			command = append(command, segment)
			segment = ""
			continue
		}

		segment += string(c)
	}

	return append(command, segment)
}
