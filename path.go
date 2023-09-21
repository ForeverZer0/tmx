package tmx

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PathResolve provides a mechanism for users to supply their own logic for resolving paths. This
// can be useful when the built-in path resolving does not satisfy the needs of the project.
//
// Because it expects a reader returned and not a path, this doubles as a way to implement a
// "virtual" filesystem, where the path need not even exist on disk, such as files being embedded
// into the binary as an internal resource.
//
// Any non-nil instance that is returned will be closed automatically by the library.
var PathResolve func(path string) (io.ReadCloser, Format, error)

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
		// If so, attempt to convert it to an absolute path if not already
		if !filepath.IsAbs(path) {
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

// DetectExt attempts to determine the format based on its file extension, falling back
// and attempting to use the file contents when locatable on disk.
func DetectExt(path string) Format {
	ext := strings.ToLower(filepath.Ext(path))

	// Determine by file extension
	switch ext {
	case "":
		return FormatUnknown
	case ".tmx", ".tsx", ".tx", ".xml":
		return FormatXML
	case ".tmj", ".tsj", ".tj", ".json":
		return FormatJSON
	}
	
	// Fallback to detecting by file contents
	if abs, err := FindPath(path); err != nil {
		return FormatUnknown
	} else if file, err := os.Open(abs); err != nil {
		return FormatUnknown
	} else {
		defer file.Close()
		return DetectContents(file)
	}
}

// DetectContents attempted to determine the format based on the text contents.
func DetectContents(reader io.ReadSeeker) Format {
	// Record the current position
	pos, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return FormatUnknown
	}
	defer reader.Seek(pos, io.SeekStart)

	// Scan characters until a '<', '{', or '[' is encountered
	r := bufio.NewReader(reader)
	for c, _, err := r.ReadRune(); err == nil; c, _, err = r.ReadRune() {
		switch c {
		case '<':
			return FormatXML
		case '{', '[':
			return FormatJSON
		}
	}

	return FormatUnknown
}

// getStream finds the given path, returning a reader object, its resolved path, and
// detected TMX format.
func getStream(abs string) (reader io.ReadCloser, ft Format, err error) {
	if reader, err = os.Open(abs); err == nil {
		ft = DetectExt(abs)
		return
	} else if PathResolve != nil {
		reader, ft, err = PathResolve(abs)
	}

	return
}

// vim: ts=4
