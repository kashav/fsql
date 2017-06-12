package transform

import (
	"encoding/hex"
	"hash"
	"io/ioutil"
	"os"
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

// computeHash applies the hash h to the file located at path. Returns a line
// of dashes for directories.
func computeHash(info os.FileInfo, path string, h hash.Hash) (interface{}, error) {
	if info.IsDir() {
		return strings.Repeat("-", h.Size()*2), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
