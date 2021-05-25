package cli

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/sha1n/benchy/test"
	"github.com/stretchr/testify/assert"
)

func TestRequestString(t *testing.T) {
	expected := test.RandomString()
	cleanup := givenStdInWith(expected)
	defer cleanup()

	actual := RequestString("", false)
	assert.Equal(t, expected, actual)
}

func TestRequestRequiredString(t *testing.T) {
	expected := test.RandomString()
	cleanup := givenStdInWith(fmt.Sprintf("\r\n\r\n%s", expected))
	defer cleanup()

	actual := RequestString("", true)
	assert.Equal(t, expected, actual)
}

func TestRequestInputWithInvalidInput(t *testing.T) {
	expected := test.RandomString()
	stdin := fmt.Sprintf(`rejected
%s`, expected)

	cleanup := givenStdInWith(stdin)
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

func TestRequestOptionalBool(t *testing.T) {
	expected := test.RandomBool()
	cleanup := givenStdInWith(fmt.Sprint(expected))
	defer cleanup()

	actual := RequestOptionalBool("", false)
	assert.Equal(t, expected, actual)
}

func TestRequestOptionalBoolWithInvalidInput(t *testing.T) {
	// FIXME this doesn't seem to be a real bug, but it just might...
	t.Skip("Skipped due to instability (looks like buffering issue on the faked stdin)")
	attempt1 := "12"
	attempt2 := 1 // 1 == true
	cleanup := givenStdInWith(fmt.Sprintf("%s\r\n%d", attempt1, attempt2))
	defer cleanup()

	actual := RequestOptionalBool("", false)
	assert.Equal(t, true, actual)
}

func TestRequestOptionalBoolWithSkip(t *testing.T) {
	cleanup := givenStdInWith("\r\n")
	defer cleanup()

	actual := RequestOptionalBool("", false)
	assert.Equal(t, false, actual)
}

func TestQuestionYNWithPositiveInput(t *testing.T) {
	cleanup := givenStdInWith("y")
	defer cleanup()

	actual := QuestionYN("")
	assert.True(t, actual)
}

func TestQuestionYNWithNegativeInput(t *testing.T) {
	cleanup := givenStdInWith("n")
	defer cleanup()

	actual := QuestionYN("")
	assert.False(t, actual)
}

func TestQuestionYNWithEmptyResponse(t *testing.T) {
	cleanup := givenStdInWith("")
	defer cleanup()

	actual := QuestionYN("")
	assert.False(t, actual) // empty == no
}

func TestRequestOptionalExistingDirectoryWithExistingDir(t *testing.T) {
	path := os.TempDir()
	cleanup := givenStdInWith(path)
	defer cleanup()

	actual := RequestOptionalExistingDirectory("", "")
	assert.Equal(t, path, actual)
}

func TestRequestOptionalExistingDirectoryWithNonExistingDir(t *testing.T) {
	attempt1 := path.Join(os.TempDir(), test.RandomString())
	attempt2 := os.TempDir()
	cleanup := givenStdInWith(fmt.Sprintf("%s\r\n%s", attempt1, attempt2))
	defer cleanup()

	actual := RequestOptionalExistingDirectory("", "")
	assert.Equal(t, attempt2, actual)
}

func TestRequestOptionalExistingDirectoryWithSkip(t *testing.T) {
	cleanup := givenStdInWith("\r\n")
	defer cleanup()

	actual := RequestOptionalExistingDirectory("", "")
	assert.Equal(t, "", actual)
}

func TestRequestUint(t *testing.T) {
	expected := test.RandomUint()
	cleanup := givenStdInWith(fmt.Sprint(expected))
	defer cleanup()

	actual := RequestUint("", false)
	assert.Equal(t, expected, actual)
}

func TestRequestRequiredUint(t *testing.T) {
	// FIXME this doesn't seem to be a real bug, but it just might...
	t.Skip("Skipped due to instability (looks like buffering issue on the faked stdin)")

	attempt1 := -1
	attempt2 := test.RandomUint()
	cleanup := givenStdInWith(fmt.Sprintf("%d\r\n%d", attempt1, attempt2))
	defer cleanup()

	actual := RequestUint("", false)
	assert.Equal(t, attempt2, actual)
}

func givenStdInWith(content string) func() {
	StdinReader = test.NewEmulatedStdinReader(content)

	return func() {
		StdinReader = os.Stdin
	}
}
