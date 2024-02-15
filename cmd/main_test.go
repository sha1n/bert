package main

import (
	"os"
	"testing"

	"github.com/sha1n/termite"
	"github.com/stretchr/testify/assert"
)

func TestExitCodeWithMissingRequiredArguments(t *testing.T) {
	expectedPanicExitCode := 1
	actualExitCode := 0

	os.Args = []string{}
	doRun(func(i int) {
		actualExitCode = i
	})

	assert.Equal(t, expectedPanicExitCode, actualExitCode)
}

func TestExitCodeWithFailFastAndCommandFailure(t *testing.T) {
	var (
		actualExitCode int
		hasValue       bool
	)

	testWith(
		t,
		[]string{
			"program",
			"-e",
			"1",
			"-k",
			"'<non-existent-command>'",
		},
		true,
		func(t *testing.T) {
			doRun(func(code int) {
				if !hasValue {
					actualExitCode = code
					hasValue = true
				}
			})
		},
	)

	assert.Equal(t, 0, actualExitCode)
}

func TestSanity(t *testing.T) {
	testWith(
		t,
		[]string{
			"program",
			"-c",
			"../test/data/integration.yaml",
			"--debug",
			"--pipe-stdout",
			"--pipe-stderr",
		},
		true,
		func(t *testing.T) {
			doRun(func(code int) {
				assert.Equal(t, 0, code)
			})
		},
	)
}

func testWith(t *testing.T, args []string, tty bool, test func(t *testing.T)) {
	origTtyValue := termite.Tty
	origOsArgs := os.Args
	termite.Tty = true

	defer func() {
		termite.Tty = origTtyValue
		os.Args = origOsArgs
	}()

	os.Args = args
	termite.Tty = tty

	test(t)
}
