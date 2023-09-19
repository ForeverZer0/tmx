package tmx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PathResolve provides a mechanism for users to supply their own logic for resolving paths. This
// can be useful when the built-in path resolving does not satisfy the needs of the project.
//
// Because it expects a reader returned and not a path, this doubles as a way to implement a
// "virtual" filesystem, where the path need not even exist on disk, such as files being embedded
// into the binary as an internal resource.
var PathResolve func(path string) io.ReadCloser

// IncludePaths contains paths to directories that will be searched when resolving relative
// file paths.
//
// More directories can be appended as needed to accomodate your project's structure.
var IncludePaths []string

// FindPath uses the given (relative or absolute) path (e.g. "../../tilesets/mountains.tsx") and
// tests whether it exists or not. Upon success, the absolute path is returned, otherwise a non-nil
// error will be present.
//
// The following locations are searched in the specified order, for both the relative path and its
// base filename:
//
//   - Any base paths passed as an argument, in the order given
//   - The current working directory
//   - Directories specified in tmx.IncludePaths
//
// os.Stat is used to test the existence of the given file/directory before returning success, so
// with minor exceptions (i.e. race conditions), it can be reasonably assumed that the path exists
// in the filesystem when no error is returned.
func FindPath(path string, base ...string) (string, error) {
	// Check if path can be resolved as-is
	if _, err := os.Stat(path); err == nil {
		if !filepath.IsAbs(path) {
			// Attempt to convert it to an absolute path if not already
			if abs, err := filepath.Abs(path); err == nil {
				return abs, nil
			}
		}
		return path, nil
	}

	if working, err := os.Getwd(); err == nil {
		base = append(base, working)
	}

	base = append(base, IncludePaths...)
	basename := filepath.Base(path)

	for _, base := range base {

		joined := filepath.Join(base, path)
		if abs, err := filepath.Abs(joined); err == nil {
			if _, err := os.Stat(abs); err == nil {
				return abs, nil
			}
		}

		joined = filepath.Join(base, basename)
		if abs, err := filepath.Abs(joined); err == nil {
			if _, err := os.Stat(abs); err == nil {
				return abs, nil
			}
		}
	}

	return path, fmt.Errorf(`%w: "%s"`, os.ErrNotExist, path)
}

// vim: ts=4
