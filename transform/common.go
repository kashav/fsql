package transform

import (
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

// truncate returns the first n letters of the string str.
// if n is greater than the size of str, return str unaltered.
// if n is < 0, return str unaltered.
func truncate(str string, n int) string {
	if len(str) < n || n < 0 {
		return str
	}

	return str[0:n]
}
