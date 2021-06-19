package pkg

import gommonsos "github.com/sha1n/gommons/pkg/os"

// ExpandUserPath attempts to expand the specified user path.
// Panics if the input is ok, but the user home resolution fails.
func ExpandUserPath(path string) (expandedPath string) {
	var err error
	if expandedPath, err = gommonsos.ExpandUserPath(path); err != nil {
		panic(err)
	}

	return
}
