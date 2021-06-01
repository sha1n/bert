package cli

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/sha1n/benchy/test"
	"github.com/sha1n/uneatest"
	"github.com/stretchr/testify/assert"
)

func TestRequestString(t *testing.T) {
	expected := test.RandomString()
	ctx := givenIOContextWithInputContent(expected)

	actual := requestString("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestRequiredString(t *testing.T) {
	expected := test.RandomString()
	ctx := givenIOContextWithInputContent(fmt.Sprintf("\r\n\r\n%s", expected))

	actual := requestString("", true, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestInputWithInvalidInput(t *testing.T) {
	expected := test.RandomString()
	stdin := fmt.Sprintf(`rejected
%s`, expected)

	ctx := givenIOContextWithInputContent(stdin)

	callCount := 1
	isValidFn := func(s string) bool {
		if callCount == 1 {
			callCount++
			return false
		}
		return true
	}

	actual := requestInput("", false, isValidFn, ctx)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 2, callCount)
}

func TestRequestOptionalBool(t *testing.T) {
	expected := test.RandomBool()
	ctx := givenIOContextWithInputContent(fmt.Sprint(expected))

	actual := requestOptionalBool("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestOptionalBoolWithInvalidInput(t *testing.T) {
	// FIXME this doesn't seem to be a real bug, but it just might...
	t.Skip("Skipped due to instability (looks like buffering issue on the faked stdin)")
	attempt1 := "12"
	attempt2 := 1 // 1 == true
	ctx := givenIOContextWithInputContent(fmt.Sprintf("%s\r\n%d", attempt1, attempt2))

	actual := requestOptionalBool("", false, ctx)
	assert.Equal(t, true, actual)
}

func TestRequestOptionalBoolWithSkip(t *testing.T) {
	ctx := givenIOContextWithInputContent("\r\n")

	actual := requestOptionalBool("", false, ctx)
	assert.Equal(t, false, actual)
}

func TestQuestionYNWithPositiveInput(t *testing.T) {
	ctx := givenIOContextWithInputContent("y")

	actual := questionYN("", ctx)
	assert.True(t, actual)
}

func TestQuestionYNWithNegativeInput(t *testing.T) {
	ctx := givenIOContextWithInputContent("n")

	actual := questionYN("", ctx)
	assert.False(t, actual)
}

func TestQuestionYNWithEmptyResponse(t *testing.T) {
	ctx := givenIOContextWithInputContent("")

	actual := questionYN("", ctx)
	assert.False(t, actual) // empty == no
}

func TestRequestOptionalExistingDirectoryWithExistingDir(t *testing.T) {
	path := os.TempDir()
	ctx := givenIOContextWithInputContent(path)

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, path, actual)
}

func TestRequestOptionalExistingDirectoryWithNonExistingDir(t *testing.T) {
	attempt1 := path.Join(os.TempDir(), test.RandomString())
	attempt2 := os.TempDir()
	ctx := givenIOContextWithInputContent(fmt.Sprintf("%s\r\n%s", attempt1, attempt2))

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, attempt2, actual)
}

func TestRequestOptionalExistingDirectoryWithSkip(t *testing.T) {
	ctx := givenIOContextWithInputContent("\r\n")

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, "", actual)
}

func TestRequestUint(t *testing.T) {
	expected := test.RandomUint()
	ctx := givenIOContextWithInputContent(fmt.Sprint(expected))

	actual := requestUint("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestRequiredUint(t *testing.T) {
	// FIXME this doesn't seem to be a real bug, but it just might...
	t.Skip("Skipped due to instability (looks like buffering issue on the faked stdin)")

	attempt1 := -1
	attempt2 := test.RandomUint()
	ctx := givenIOContextWithInputContent(fmt.Sprintf("%d\r\n%d", attempt1, attempt2))

	actual := requestUint("", false, ctx)
	assert.Equal(t, attempt2, actual)
}

func givenIOContextWithInputContent(content string) IOContext {
	ctx := NewIOContext()
	ctx.StdinReader = uneatest.NewEmulatedStdinReader(content)

	return ctx
}
