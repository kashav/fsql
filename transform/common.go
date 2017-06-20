package transform

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// formatName runs the correct name format function based on the value of arg.
func formatName(arg, name string) interface{} {
	switch strings.ToUpper(arg) {
	case "UPPER":
		return upper(name)
	case "LOWER":
		return lower(name)
	}
	return nil
}

// upper returns the uppercased version of name.
func upper(name string) interface{} {
	return strings.ToUpper(name)
}

// lower returns the lowercase version of name.
func lower(name string) interface{} {
	return strings.ToLower(name)
}

// truncate returns the first n characters of str. If n is greater than the
// length of str or less than 0, return str.
func truncate(str string, n int) string {
	if len(str) < n || n < 0 {
		return str
	}

	return str[0:n]
}

// FindHash returns a func to create a new hash based on the provided name.
func FindHash(name string) func() hash.Hash {
	switch strings.ToUpper(name) {
	case "SHA1":
		return sha1.New
	}
	return nil
}

// ComputeHash applies the hash h to the file located at path. Returns a line
// of dashes for directories.
func ComputeHash(info os.FileInfo, path string, h hash.Hash) (interface{}, error) {
	fallback := strings.Repeat("-", h.Size()*2)

	// If the current file is a symlink, attempt to evaluate the link and
	// stat the resultant file. If either process fails, ignore the error and
	// return the fallback.
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		var err error
		if path, err = filepath.EvalSymlinks(path); err != nil {
			return fallback, nil
		}
		if info, err = os.Stat(path); err != nil {
			return fallback, nil
		}
	}

	if info.IsDir() {
		return fallback, nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if _, err := h.Write(b); err != nil {
		return nil, err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
