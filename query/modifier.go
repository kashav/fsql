package query

import (
	"fmt"
	"os"
	_ "path/filepath"
	"strings"
	"time"
)

type Modifier struct {
	Name      string
	Arguments []string
}

func (m *Modifier) String() string {
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(m.Arguments, ", "))
}

type modifierParams struct {
	Key   string
	Input interface{}
	Path  string
	Info  os.FileInfo
	Args  []string
}

func formatSize(p *modifierParams) (interface{}, error) {
	size, _ := p.Input.(int64)
	var s string
	switch p.Args[0] {
	case "kb":
		s = fmt.Sprintf("%fkb", float64(size)/(1<<10))
	case "mb":
		s = fmt.Sprintf("%fmb", float64(size)/(1<<20))
	case "gb":
		s = fmt.Sprintf("%fgb", float64(size)/(1<<30))
	default:
		s = string(size)
	}
	return s, nil
}

func formatTime(p *modifierParams) (interface{}, error) {
	var t string
	switch p.Args[0] {
	case "iso":
		t = p.Info.ModTime().Format(time.RFC3339)
	case "unix":
		t = p.Info.ModTime().Format(time.UnixDate)
	default:
		t = p.Info.ModTime().Format(time.Stamp)
	}
	return t, nil
}

func format(p *modifierParams) (interface{}, error) {
	switch p.Key {
	case "size":
		return formatSize(p)
	case "time":
		return formatTime(p)
	default:
		return nil,
			fmt.Errorf("Function FORMAT not implemented for attribute %s.\n",
				p.Key)
	}
}

func upper(p *modifierParams) (interface{}, error) {
	if p.Key != "name" {
		return nil, fmt.Errorf("Function UPPER not implemented for attribute %s.\n",
			p.Key)
	}

	return strings.ToUpper(p.Input.(string)), nil
}

func fullpath(p *modifierParams) (interface{}, error) {
	if p.Key != "name" {
		return nil,
			fmt.Errorf("Function FULLPATH not implemented for attribute %s.\n",
				p.Key)
	}

	return p.Path, nil
}

func defaultValue(key, path string, info os.FileInfo) interface{} {
	switch key {
	case "mode":
		return info.Mode()
	case "size":
		return info.Size()
	case "time":
		return info.ModTime().Format(time.Stamp)
	case "name":
		return info.Name()
	}

	return nil
}

var modifierFns = map[string]func(*modifierParams) (interface{}, error){
	"format":   format,
	"upper":    upper,
	"fullpath": fullpath,
}

// ApplyModifiers iterates through each SELECT attribute for this query
// and applies the associated modifier to the attribute's value.
func (q *Query) ApplyModifiers(path string, info os.FileInfo) map[string]interface{} {
	results := make(map[string]interface{}, len(q.Attributes))

	for k, _ := range q.Attributes {
		value := defaultValue(k, path, info)

		modifiers, ok := q.Modifiers[k]
		if !ok {
			results[k] = value
			continue
		}

		for _, m := range modifiers {
			fn, ok := modifierFns[m.Name]
			if !ok {
				fmt.Fprintf(os.Stderr, "Function %s is not implemented.\n",
					strings.ToUpper(m.Name))
				os.Exit(1)
			}

			var err error
			value, err = fn(&modifierParams{k, value, path, info, m.Arguments})
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		}

		results[k] = value
	}

	return results
}