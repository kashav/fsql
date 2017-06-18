package transform

import (
	"hash"
	"os"
	"reflect"
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
// Depending on the type of p.Value, we may recursively run this method
// on every element of the structure.
//
// We're using reflect _quite_ heavily for this, meaning it's kind of unsafe,
// it'd be great if we could find another solution while keeping it as
// abstract as it is.
func Parse(p *ParseParams) (val interface{}, err error) {
	kind := reflect.TypeOf(p.Value).Kind()

	// If we have a slice/array, recursively run Parse on each element.
	if kind == reflect.Slice || kind == reflect.Array {
		s := reflect.ValueOf(p.Value)
		for i := 0; i < s.Len(); i++ {
			p.Value = s.Index(i).Interface()
			if val, err = Parse(p); err != nil {
				return nil, err
			}
			s.Index(i).Set(reflect.ValueOf(val))
		}
		return s.Interface(), nil
	}

	// If we have a map, recursively run Parse on each *key* and create a new
	// map out of the return values.
	if kind == reflect.Map {
		result := reflect.MakeMap(reflect.TypeOf(p.Value))
		for _, key := range reflect.ValueOf(p.Value).MapKeys() {
			p.Value = key.Interface()
			if val, err = Parse(p); err != nil {
				return nil, err
			}
			result.SetMapIndex(reflect.ValueOf(val), reflect.ValueOf(true))
		}
		return result.Interface(), nil
	}

	// Not a slice nor a map.
	switch strings.ToUpper(p.Name) {
	case "FORMAT":
		val, err = p.format()
	case "UPPER":
		val = upper(p.Value.(string))
	case "LOWER":
		val = lower(p.Value.(string))
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

// format runs the correct format function based on the provided attribute.
func (p *ParseParams) format() (val interface{}, err error) {
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

// formatSize formats the size attribute. Valid arguments include `B`, `KB`,
// `MB`, and `GB` (case insensitive).
func (p *ParseParams) formatSize() (interface{}, error) {
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
		size *= 1 << 30
	default:
		return nil, nil
	}
	return size, nil
}

// formatTime formats the time attribute. Valid arguments include `ISO`,
// `UNIX`, (case insensitive) or a custom layout. If a custom layout is
// provided, it must be set according to 2006-01-02T15:04:05.999999-07:00.
func (p *ParseParams) formatTime() (interface{}, error) {
	var (
		t   time.Time
		err error
	)
	switch strings.ToUpper(p.Args[0]) {
	case "ISO":
		t, err = time.Parse(time.RFC3339, p.Value.(string))
	case "UNIX":
		t, err = time.Parse(time.UnixDate, p.Value.(string))
	default:
		t, err = time.Parse(p.Args[0], p.Value.(string))
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *ParseParams) hash(h hash.Hash) (interface{}, error) {
	info, err := os.Stat(p.Value.(string))
	if err != nil {
		return nil, err
	}
	return ComputeHash(info, p.Value.(string), h)
}
