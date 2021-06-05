package main

import (
	"os"
	"testing"

	"github.com/sha1n/termite"
	"github.com/stretchr/testify/assert"
)

func TestExitCodeWithMissingRequiredArguments(t *testing.T) {
	expectedPanicExitCode := 1
	actualExitcode := 0

	os.Args = []string{}
	doRun(func(i int) {
		actualExitcode = i
	})

	assert.Equal(t, expectedPanicExitCode, actualExitcode)
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
