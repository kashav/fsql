package transform

import (
	"strings"
)

func formatName(arg, name string) string {
	switch strings.ToUpper(arg) {
	case "UPPER":
		return upper(name)
	case "LOWER":
		return lower(name)
	}
	return ""
}

func upper(name string) string {
	return strings.ToUpper(name)
}

func lower(name string) string {
	return strings.ToLower(name)
}
