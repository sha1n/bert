package cli

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sha1n/benchy/internal/github"
	"github.com/sha1n/benchy/test"
	"github.com/stretchr/testify/assert"
)

const (
	v1_0_0 = "v1.0.0"
	v1_0_1 = "v1.0.1"
)

func TestRunSelfUpdateWithReleaseError(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		expectedError := errors.New(test.RandomString())
		getLatestRelease := aGetLatestReleaseFnWith(nil, expectedError)

		actualError := runSelfUpdateWith(test.RandomString(), test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.Equal(t, expectedError, actualError)
	})
}

func TestRunSelfUpdateWithCurrentRelease(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_0
		latestRelease := &fakeRelease{tag: v1_0_0}
		getLatestRelease := aGetLatestReleaseFnWith(latestRelease, nil)

		actualError := runSelfUpdateWith(currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.False(t, latestRelease.downloadCalled)
	})
}

func TestRunSelfUpdateWithDownloadError(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		currentVersion := v1_0_0
		expectedError := errors.New(test.RandomString())
		latestRelease := &fakeRelease{tag: v1_0_1, downloadError: expectedError}
		getLatestRelease := aGetLatestReleaseFnWith(latestRelease, nil)

		actualError := runSelfUpdateWith(currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.Error(t, actualError)
		assert.Equal(t, expectedError, actualError)
		assert.True(t, latestRelease.downloadCalled)
	})
}

func TestRunSelfUpdateWithSuccessfulDownload(t *testing.T) {
	runWith(func(resolveBinaryPathFn ResolveBinaryPathFn) {
		expectedFileContent := []byte(test.RandomString())
		currentVersion := v1_0_0
		latestRelease := &fakeRelease{tag: v1_0_1, data: expectedFileContent}
		getLatestRelease := aGetLatestReleaseFnWith(latestRelease, nil)

		actualError := runSelfUpdateWith(currentVersion, test.RandomString(), resolveBinaryPathFn, getLatestRelease)

		assert.NoError(t, actualError)
		assert.True(t, latestRelease.downloadCalled)

		path, err := resolveBinaryPathFn()
		assert.NoError(t, err)

		actualFileContent, err := os.ReadFile(path)
		assert.NoError(t, err)

		assert.Equal(t, expectedFileContent, actualFileContent)
	})
}

func aGetLatestReleaseFnWith(r github.Release, e error) GetLatestReleaseFn {
	return func() (github.Release, error) {
		return r, e
	}
}

func runWith(doTest func(resolveBinaryPathFn ResolveBinaryPathFn)) {
	resolveBinaryPathFn, cleanup := resolveExecutableFn()
	defer cleanup()

	doTest(resolveBinaryPathFn)
}

type fakeRelease struct {
	tag            string
	downloadError  error
	downloadCalled bool
	data           []byte
}

func (r *fakeRelease) TagName() string {
	return r.tag
}

func (r *fakeRelease) DownloadAsset() (rc io.ReadCloser, err error) {
	r.downloadCalled = true
	rc, err = nil, r.downloadError

	if r.downloadError == nil {
		rc = ioutil.NopCloser(bytes.NewReader(r.data))
	}

	return rc, err
}

func resolveExecutableFn() (ResolveBinaryPathFn, func()) {
	f, _ := ioutil.TempFile("", "fake_binary")

	fn := func() (string, error) {
		return f.Name(), nil
	}

	return fn, func() { os.Remove(f.Name()) }
}
