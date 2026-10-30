package schema

import (
	"errors"
	"io/fs"
	"os"
)

// exists reports whether the file exists.
func exists(p string) (bool, error) {
	if _, err := os.Stat(p); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err //nolint:wrapcheck
	}
	return true, nil
}
