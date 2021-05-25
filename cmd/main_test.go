package main

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var runMainCommand = []string{
	"go",
	"run",
	"-mod=readonly",
	"main.go",
}

func TestExitCodeWhenRequiredConfigArgIsMissing(t *testing.T) {
	expectedExitCode := 1
	buf := new(bytes.Buffer)	

	cmd := exec.Command(runMainCommand[0], runMainCommand[1:]...)

	cmd.Stdout = buf
	cmd.Stderr = buf

	assert.NoError(t, cmd.Start())
	state, err := cmd.Process.Wait()

	assert.Contains(t, buf.String(), "Error: required flag(s) \"config\" not set")
	assert.NoError(t, err)
	assert.Equal(t, expectedExitCode, state.ExitCode())

}
