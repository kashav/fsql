package transform

import (
	"crypto"
	_ "crypto/sha1" //Import SHA-1 hashing function
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
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

func (p FormatParams) argAs(index int, kind reflect.Kind) (interface{}, error) {
	// return args as an int
	asInt := func() (interface{}, error) {
		var n int
		var err error
		if err == nil && len(p.Args) > 0 && p.Args[index] != "" {
			if n, err = strconv.Atoi(p.Args[index]); err != nil {
				return nil, err
			}
		}
		return n, nil
	}
	switch kind {
	case reflect.Int:
		return asInt()
	default:
		return nil, &ErrNotImplemented{p.Name, p.Attribute}
	}
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
		val, err = p.hash(crypto.SHA1)
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

// hash will take the hash function, based on the hasher type supplied.
func (p *FormatParams) hash(hasher crypto.SignerOpts) (interface{}, error) {
	var err error
	var n interface{}
	var h interface{}

	n, err = p.argAs(0, reflect.Int)
	if err != nil {
		return nil, err
	}
	h, err = hash(p.Info, p.Path, hasher)
	return truncate(h.(string), n.(int)), nil
}

func hash(info os.FileInfo, path string, hasher crypto.SignerOpts) (interface{}, error) {
	if info.IsDir() {
		return strings.Repeat("-", hasher.HashFunc().Size()), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	h := hasher.HashFunc().New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// DefaultFormatValue returns the default format value for the provided
// attribute attr based on path and info.
func DefaultFormatValue(attr, path string, info os.FileInfo) (interface{}, error) {
	switch attr {
	case "mode":
		return info.Mode(), nil
	case "name":
		return info.Name(), nil
	case "size":
		return info.Size(), nil
	case "time":
		return info.ModTime().Format(time.Stamp), nil
	case "hash":
		v, err := hash(info, path, crypto.SHA1)
		if err == nil {
			return truncate(v.(string), 7), nil
		}
	}
	return nil, &ErrUnsupportedFormat{"", attr}
}
