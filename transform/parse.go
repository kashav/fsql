package transform

import (
	"strconv"
	"strings"
	"time"
)

// ParseParams holds the params for a parse-modifier function.
type ParseParams struct {
	Attribute string
	Value     interface{}

	Name string
	Args []string
}

// Parse runs the associated modifier function for the provided parameters.
func Parse(p *ParseParams) (val interface{}, err error) {
	switch strings.ToUpper(p.Name) {
	case "FORMAT":
		val, err = pFormat(p)
	case "UPPER":
		val, err = upper(p.Value.(string)), nil
	case "LOWER":
		val, err = lower(p.Value.(string)), nil
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrNotImplemented{p.Name, p.Attribute}
	}
	return val, nil
}

func pFormat(p *ParseParams) (val interface{}, err error) {
	switch p.Attribute {
	case "name":
		val, err = formatName(p.Args[0], p.Value.(string)), nil
	case "size":
		val, err = pFormatSize(p)
	case "time":
		val, err = pFormatTime(p)
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrUnsupportedFormat{p.Args[0], p.Attribute}
	}
	return val, nil
}

func pFormatSize(p *ParseParams) (interface{}, error) {
	size, err := strconv.ParseFloat(p.Value.(string), 64)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(p.Args[0]) {
	case "B":
		size *= 1
	case "KB":
		size *= 1 << 10
	case "MB":
		size *= 1 << 20
	case "GB":
		size *= 1 << 20
	default:
		return nil, nil
	}

	return size, nil
}

func pFormatTime(p *ParseParams) (interface{}, error) {
	var t time.Time
	var err error

	switch strings.ToUpper(p.Args[0]) {
	case "ISO":
		t, err = time.Parse(time.RFC3339, p.Value.(string))
	case "UNIX":
		t, err = time.Parse(time.UnixDate, p.Value.(string))
	default:
		t, err = time.Parse("Jan 02 2006 15 04", p.Value.(string))
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}
