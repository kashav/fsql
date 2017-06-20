package transform

import (
	"fmt"
	"hash"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultHashLength = 7

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
	case "SHA1":
		val, err = p.hash(FindHash(p.Name)())
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
// insensitive), or a custom layout layout. If a custom layout is provided, it
// must be set according to 2006-01-02T15:04:05.999999-07:00.
func (p *FormatParams) formatTime() (interface{}, error) {
	switch strings.ToUpper(p.Args[0]) {
	case "ISO":
		return p.Info.ModTime().Format(time.RFC3339), nil
	case "UNIX":
		return p.Info.ModTime().Format(time.UnixDate), nil
	default:
		return p.Info.ModTime().Format(p.Args[0]), nil
	}
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

// hash applies the provided hash algorithm h with ComputeHash.
func (p *FormatParams) hash(h hash.Hash) (interface{}, error) {
	var (
		err    error
		n      int
		result interface{}
	)

	if len(p.Args) == 0 || p.Args[0] == "" {
		n = defaultHashLength
	} else if strings.ToUpper(p.Args[0]) == "FULL" {
		n = -1
	} else if n, err = strconv.Atoi(p.Args[0]); err != nil {
		return nil, err
	}

	if result, err = ComputeHash(p.Info, p.Path, h); err != nil {
		return nil, err
	}

	return truncate(result.(string), n), nil
}

// DefaultFormatValue returns the default format value for the provided
// attribute attr based on path and info.
func DefaultFormatValue(attr, path string, info os.FileInfo) (value interface{}, err error) {
	switch attr {
	case "mode":
		value = info.Mode()
	case "name":
		value = info.Name()
	case "size":
		value = info.Size()
	case "time":
		value = info.ModTime().Format(time.Stamp)
	case "hash":
		if value, err = ComputeHash(info, path, FindHash("SHA1")()); value != nil {
			value = truncate(value.(string), defaultHashLength)
		}
	default:
		err = fmt.Errorf("unknown attribute %s", attr)
	}
	return value, err
}
