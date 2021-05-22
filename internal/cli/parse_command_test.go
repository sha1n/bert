package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	command := "do something"
	expected := []string{
		"do",
		"something",
	}
	assert.Equal(t, expected, parseCommand(command))
}

func TestParseCommandWithSingleQuotes(t *testing.T) {
	command := "do 'whatever is required'"
	expected := []string{
		"do",
		"whatever is required",
	}
	assert.Equal(t, expected, parseCommand(command))
}

func TestParseCommandWithDoubleQuotes(t *testing.T) {
	command := "do \"whatever is required\""
	expected := []string{
		"do",
		"whatever is required",
	}
	assert.Equal(t, expected, parseCommand(command))
}

func TestParseCommandWithMixedQuotes(t *testing.T) {
	command := "do \"someone else's work\""
	expected := []string{
		"do",
		"someone else's work",
	}
	assert.Equal(t, expected, parseCommand(command))
}

func TestParseCommandWithEscapedDoubleQuotes(t *testing.T) {
	command := "c \"a \\\" b\""
	expected := []string{
		"c",
		"a \" b",
	}
	assert.Equal(t, expected, parseCommand(command))
}

func TestParseCommandWithEscapedQuotes(t *testing.T) {
	command := "c 'a \\' b'"
	expected := []string{
		"c",
		"a ' b",
	}
	assert.Equal(t, expected, parseCommand(command))
}
