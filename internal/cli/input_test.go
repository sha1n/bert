package cli

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/sha1n/benchy/api"
	clibtest "github.com/sha1n/clib/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRequestString(t *testing.T) {
	expected := clibtest.RandomString()
	ctx := givenIOContextWithInputContent(expected)

	actual := requestString("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestRequiredString(t *testing.T) {
	expected := clibtest.RandomString()
	ctx := givenIOContextWithInputContent(fmt.Sprintf("\r\n\r\n%s", expected))

	actual := requestString("", true, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestInputWithInvalidInput(t *testing.T) {
	expected := clibtest.RandomString()
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
	expected := clibtest.RandomBool()
	ctx := givenIOContextWithInputContent(fmt.Sprint(expected))

	actual := requestOptionalBool("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestOptionalBoolWithInvalidInput(t *testing.T) {
	attempt1 := "12"
	attempt2 := 1 // 1 == true
	userInputSequence := fmt.Sprintf(`%s
%d`, attempt1, attempt2)
	ctx := givenIOContextWithInputContent(userInputSequence)

	actual := requestOptionalBool("", false, ctx)
	assert.True(t, actual)
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
	attempt1 := path.Join(os.TempDir(), clibtest.RandomString())
	attempt2 := os.TempDir()
	userInputSequence := fmt.Sprintf(`%s
n
%s`, attempt1, attempt2)

	ctx := givenIOContextWithInputContent(userInputSequence)

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, attempt2, actual)
}

func TestRequestOptionalExistingDirectoryWithNonExistingDirAndAutoCreation(t *testing.T) {
	userEnteredNonExistingDir := path.Join(os.TempDir(), clibtest.RandomString())
	userInputSequence := fmt.Sprintf(`%s
y`,
		userEnteredNonExistingDir,
	)

	ctx := givenIOContextWithInputContent(userInputSequence)

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, userEnteredNonExistingDir, actual)
}

func TestRequestOptionalExistingDirectoryWithSkip(t *testing.T) {
	ctx := givenIOContextWithInputContent("\r\n")

	actual := requestOptionalExistingDirectory("", "", ctx)
	assert.Equal(t, "", actual)
}

func TestRequestUint(t *testing.T) {
	expected := clibtest.RandomUint()
	ctx := givenIOContextWithInputContent(fmt.Sprint(expected))

	actual := requestUint("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestUintWithEmptyInput(t *testing.T) {
	expected := uint(0)
	ctx := givenIOContextWithInputContent("\r\n")

	actual := requestUint("", false, ctx)
	assert.Equal(t, expected, actual)
}

func TestRequestUintWithInvalidInput(t *testing.T) {
	expected := uint(2)
	userInputSequence := fmt.Sprintf(`invalid
-1
%d`,
		expected,
	)

	ctx := givenIOContextWithInputContent(userInputSequence)

	actual := requestUint("", true, ctx)
	assert.Equal(t, expected, actual)
}

func givenIOContextWithInputContent(content string) api.IOContext {
	ctx := api.NewIOContext()
	ctx.StdinReader = clibtest.NewEmulatedStdinReader(content)

	return ctx
}
