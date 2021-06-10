package pkg

import clibos "github.com/sha1n/clib/pkg/os"

// ExpandUserPath attempts to expand the specified user path.
// Panics if the input is ok, but the user home resolution fails.
func ExpandUserPath(path string) (expandedPath string) {
	var err error
	if expandedPath, err = clibos.ExpandUserPath(path); err != nil {
		panic(err)
	}

	return
}
