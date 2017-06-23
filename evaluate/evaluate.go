package evaluate

import (
	"os"
	"strconv"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

// Opts represents a set of options used in the evaluate functions.
type Opts struct {
	Path      string
	File      os.FileInfo
	Attribute string
	Modifiers []Modifier
	Operator  tokenizer.TokenType
	Value     interface{}
}

// Modifier represents an attribute modifier.
type Modifier struct {
	Name      string
	Arguments []string
}

// Evaluate runs the respective evaluate function for the provided options.
func Evaluate(o *Opts) (bool, error) {
	switch o.Attribute {
	case "name":
		return evaluateName(o)
	case "size":
		return evaluateSize(o)
	case "time":
		return evaluateTime(o)
	case "mode":
		return evaluateMode(o)
	case "hash":
		return evaluateHash(o)
	}
	return false, &ErrUnsupportedAttribute{o.Attribute}
}

// evaluateName evaluates a Condition with attribute `name`.
func evaluateName(o *Opts) (bool, error) {
	var a, b interface{}
	switch o.Value.(type) {
	case string, []string, map[interface{}]bool:
		a = o.File.Name()
		b = o.Value
	default:
		return false, &ErrUnsupportedType{o.Attribute, o.Value}
	}
	return cmpAlpha(o, a, b)
}

// evaluateSize evaluates a Condition with attribute `size`.
func evaluateSize(o *Opts) (bool, error) {
	var a, b interface{}
	switch o.Value.(type) {
	case float64:
		a = o.File.Size()
		b = int64(o.Value.(float64))
	case map[interface{}]bool:
		a = o.File.Size()
		b = o.Value
	case string:
		size, err := strconv.ParseFloat(o.Value.(string), 10)
		if err != nil {
			return false, err
		}
		a = o.File.Size()
		b = int64(size)
	default:
		return false, &ErrUnsupportedType{o.Attribute, o.Value}
	}
	return cmpNumeric(o, a, b)
}

// evaluateTime evaluates a Condition with attribute `time`.
func evaluateTime(o *Opts) (bool, error) {
	var a, b interface{}
	switch o.Value.(type) {
	case string:
		t, err := time.Parse("Jan 02 2006 15 04", o.Value.(string))
		if err != nil {
			return false, err
		}
		a = o.File.ModTime()
		b = t
	case map[interface{}]bool, time.Time:
		a = o.File.ModTime()
		b = o.Value
	default:
		return false, &ErrUnsupportedType{o.Attribute, o.Value}
	}
	return cmpTime(o, a, b)
}

// evaluateMode evaluates a Condition with attribute `mode`.
func evaluateMode(o *Opts) (bool, error) { return cmpMode(o) }

// evaluateHash evaluates a Condition with attribute `hash`.
func evaluateHash(o *Opts) (bool, error) { return cmpHash(o) }
