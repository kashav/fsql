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
		val, err = p.format()
	case "UPPER":
		val = upper(p.Value.(string))
	case "LOWER":
		val = lower(p.Value.(string))
	case "FULLPATH":
		val, err = p.fullPath()
	case "SHORTPATH":
		val, err = p.shortPath()
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrNotImplemented{p.Name, p.Attribute}
	}
	return val, nil
}

// format runs a format function based on the value of the provided attribute.
func (p *FormatParams) format() (val interface{}, err error) {
	switch p.Attribute {
	case "name":
		val = formatName(p.Args[0], p.Value.(string))
	case "size":
		val, err = p.formatSize()
	case "time":
		val, err = p.formatTime()
	}

	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, &ErrUnsupportedFormat{p.Args[0], p.Attribute}
	}
	return val, nil
}

// formatSize formats a size. Valid arguments include `KB`, `MB`, `GB` (case
// insensitive).
func (p *FormatParams) formatSize() (interface{}, error) {
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

// formatTime formats a time. Valid arguments include `UNIX` and `ISO` (case
// insensitive).
func (p *FormatParams) formatTime() (interface{}, error) {
	switch strings.ToUpper(p.Args[0]) {
	case "ISO":
		return p.Info.ModTime().Format(time.RFC3339), nil
	case "UNIX":
		return p.Info.ModTime().Format(time.UnixDate), nil
	}
	return nil, nil
}

// fullPath returns the full path of the current file. Only supports the
// `name` attribute.
func (p *FormatParams) fullPath() (interface{}, error) {
	if p.Attribute != "name" {
		return nil, nil
	}
	return p.Path, nil
}

// shortPath returns the short path of the current file. Only supports the
// `name` attribute.
func (p *FormatParams) shortPath() (interface{}, error) {
	if p.Attribute != "name" {
		return nil, nil
	}
	return p.Info.Name(), nil
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
