package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sha1n/benchy/test"
	"github.com/stretchr/testify/assert"
)

func TestRequestInput(t *testing.T) {
	expected := test.RandomString()
	cleanup := givenStdInWith(t, expected)
	defer cleanup()

	actual := RequestInput("", false, defaultIsValidFn)
	assert.Equal(t, expected, actual)
}

func TestRequestInputWithInvalidInput(t *testing.T) {
	expected := test.RandomString()
	stdin := fmt.Sprintf(`rejected
%s`, expected)

	cleanup := givenStdInWith(t, stdin)
	defer cleanup()

	callCount := 1
	isValidFn := func(s string) bool {
		if callCount == 1 {
			callCount++
			return false
		}
		return true
	}

	actual := RequestInput("", false, isValidFn)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 2, callCount)
}

func TestRequestBool(t *testing.T) {
	expected := test.RandomBool()
	cleanup := givenStdInWith(t, fmt.Sprint(expected))
	defer cleanup()

	actual := RequestOptionalBool("", false)
	assert.Equal(t, expected, actual)
}

func TestQuestionYNWithPositiveInput(t *testing.T) {
	cleanup := givenStdInWith(t, "y")
	defer cleanup()

	actual := QuestionYN("")
	assert.True(t, actual)
}

func TestQuestionYNWithNegativeInput(t *testing.T) {
	cleanup := givenStdInWith(t, "n")
	defer cleanup()

	actual := QuestionYN("")
	assert.False(t, actual)
}

func TestQuestionYNWithEmptyResponse(t *testing.T) {
	cleanup := givenStdInWith(t, "")
	defer cleanup()

	actual := QuestionYN("")
	assert.False(t, actual) // empty == no
}

func TestRequestUint(t *testing.T) {
	expected := test.RandomUint()
	cleanup := givenStdInWith(t, fmt.Sprint(expected))
	defer cleanup()

	actual := RequestUint("", false)
	assert.Equal(t, expected, actual)
}

func givenStdInWith(t *testing.T, content string) func() {
	var err error
	var tmpfile *os.File
	tmpfile, err = ioutil.TempFile("", "emulated_stdin")
	assert.NoError(t, err)

	_, err = tmpfile.WriteString(content)
	assert.NoError(t, err)

	_, err = tmpfile.Seek(0, 0)
	assert.NoError(t, err)

	oldStdin := os.Stdin

	os.Stdin = tmpfile

	return func() {
		os.Stdin = oldStdin
		_ = tmpfile.Close()
		os.Remove(tmpfile.Name())

	}
}
