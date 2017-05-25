package transform

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// FormatParams holds the params for a format-modifier function.
type FormatParams struct {
	Attribute string
	Path      string
	Info      os.FileInfo
	Value     interface{}

	Name string
	Args []string
}

// Format runs the respective format function on the provided parameters.
func Format(p *FormatParams) (val interface{}, err error) {
	switch strings.ToUpper(p.Name) {
	case "FORMAT":
		val, err = fFormat(p)
	case "UPPER":
		val, err = upper(p.Value.(string)), nil
	case "LOWER":
		val, err = lower(p.Value.(string)), nil
	case "FULLPATH":
		val, err = fFullPath(p)
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrNotImplemented{p.Name, p.Attribute}
	}
	return val, nil
}

func fFormat(p *FormatParams) (val interface{}, err error) {
	switch p.Attribute {
	case "name":
		val, err = formatName(p.Args[0], p.Value.(string)), nil
	case "size":
		val, err = fFormatSize(p)
	case "time":
		val, err = fFormatTime(p)
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrUnsupportedFormat{p.Args[0], p.Attribute}
	}
	return val, nil
}

func fFormatSize(p *FormatParams) (interface{}, error) {
	size := p.Value.(int64)
	switch strings.ToUpper(p.Args[0]) {
	case "KB":
		return fmt.Sprintf("%fkb", float64(size)/(1<<10)), nil
	case "MB":
		return fmt.Sprintf("%fmb", float64(size)/(1<<20)), nil
	case "GB":
		return fmt.Sprintf("%fgb", float64(size)/(1<<30)), nil
	}
	return nil, nil
}

func fFormatTime(p *FormatParams) (interface{}, error) {
	switch strings.ToUpper(p.Args[0]) {
	case "ISO":
		return p.Info.ModTime().Format(time.RFC3339), nil
	case "UNIX":
		return p.Info.ModTime().Format(time.UnixDate), nil
	}
	return nil, nil
}

func fFullPath(p *FormatParams) (interface{}, error) {
	if p.Attribute != "name" {
		return nil,
			fmt.Errorf("function FULLPATH not implemented for attribute %s",
				p.Attribute)
	}
	return p.Path, nil
}

// DefaultFormatValue returns the default format value for the provided
// attribute attr based on path and info.
func DefaultFormatValue(attr, path string, info os.FileInfo) interface{} {
	switch attr {
	case "mode":
		return info.Mode()
	case "name":
		return info.Name()
	case "size":
		return info.Size()
	case "time":
		return info.ModTime().Format(time.Stamp)
	}
	return nil
}
