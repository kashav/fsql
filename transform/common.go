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
