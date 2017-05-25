package transform

import (
	"strings"
)

func formatName(arg, name string) interface{} {
	switch strings.ToUpper(arg) {
	case "UPPER":
		return upper(name)
	case "LOWER":
		return lower(name)
	}
	return nil
}

func upper(name string) interface{} {
	return strings.ToUpper(name)
}

func lower(name string) interface{} {
	return strings.ToLower(name)
}
