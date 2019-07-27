package evaluate

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kashav/fsql/tokenizer"
	"github.com/kashav/fsql/transform"
)

// cmpAlpha performs alphabetic comparison on a and b.
func cmpAlpha(o *Opts, a, b interface{}) (result bool, err error) {
	switch o.Operator {
	case tokenizer.Equals:
		result = a.(string) == b.(string)
	case tokenizer.NotEquals:
		result = a.(string) != b.(string)
	case tokenizer.Like:
		aStr, bStr := a.(string), b.(string)
		if strings.HasPrefix(bStr, "%") && strings.HasSuffix(bStr, "%") {
			result = strings.Contains(aStr, bStr[1:len(bStr)-1])
		} else if strings.HasPrefix(bStr, "%") {
			result = strings.HasSuffix(aStr, bStr[1:])
		} else if strings.HasSuffix(bStr, "%") {
			result = strings.HasPrefix(aStr, bStr[:len(bStr)-1])
		} else {
			result = strings.Contains(aStr, bStr)
		}
	case tokenizer.RLike:
		result = regexp.MustCompile(b.(string)).MatchString(a.(string))
	case tokenizer.In:
		switch t := b.(type) {
		case map[interface{}]bool:
			if _, ok := t[a.(string)]; ok {
				result = true
			}
		case []string:
			for _, el := range t {
				if a.(string) == el {
					result = true
				}
			}
		case string:
			for _, el := range strings.Split(t, ",") {
				if a.(string) == el {
					result = true
				}
			}
		}
	default:
		err = &ErrUnsupportedOperator{o.Attribute, o.Operator}
	}
	return result, err
}

// cmpNumeric performs numeric comparison on a and b.
func cmpNumeric(o *Opts, a, b interface{}) (result bool, err error) {
	switch o.Operator {
	case tokenizer.Equals:
		result = a.(int64) == b.(int64)
	case tokenizer.NotEquals:
		result = a.(int64) != b.(int64)
	case tokenizer.GreaterThanEquals:
		result = a.(int64) >= b.(int64)
	case tokenizer.GreaterThan:
		result = a.(int64) > b.(int64)
	case tokenizer.LessThanEquals:
		result = a.(int64) <= b.(int64)
	case tokenizer.LessThan:
		result = a.(int64) < b.(int64)
	case tokenizer.In:
		if _, ok := b.(map[interface{}]bool)[a.(int64)]; ok {
			result = true
		}
	default:
		err = &ErrUnsupportedOperator{o.Attribute, o.Operator}
	}
	return result, err
}

// cmpTime performs time comparison on a and b.
func cmpTime(o *Opts, a, b interface{}) (result bool, err error) {
	switch o.Operator {
	case tokenizer.Equals:
		result = a.(time.Time).Equal(b.(time.Time))
	case tokenizer.NotEquals:
		result = !a.(time.Time).Equal(b.(time.Time))
	case tokenizer.GreaterThanEquals:
		result = a.(time.Time).After(b.(time.Time)) || a.(time.Time).Equal(b.(time.Time))
	case tokenizer.GreaterThan:
		result = a.(time.Time).After(b.(time.Time))
	case tokenizer.LessThanEquals:
		result = a.(time.Time).Before(b.(time.Time)) || a.(time.Time).Equal(b.(time.Time))
	case tokenizer.LessThan:
		result = a.(time.Time).Before(b.(time.Time))
	case tokenizer.In:
		if _, ok := b.(map[interface{}]bool)[a.(time.Time)]; ok {
			result = true
		}
	default:
		err = &ErrUnsupportedOperator{o.Attribute, o.Operator}
	}
	return result, err
}

// cmpMode performs mode comparison with info and typ.
func cmpMode(o *Opts) (result bool, err error) {
	if o.Operator != tokenizer.Is {
		return false, &ErrUnsupportedOperator{o.Attribute, o.Operator}
	}
	switch strings.ToUpper(o.Value.(string)) {
	case "DIR":
		result = o.File.Mode().IsDir()
	case "REG":
		result = o.File.Mode().IsRegular()
	default:
		result = false
	}
	return result, err
}

// cmpHash computes the hash of the current file and compares it with the
// provided value.
func cmpHash(o *Opts) (result bool, err error) {
	hashType := "SHA1"
	if len(o.Modifiers) > 0 {
		hashType = o.Modifiers[0].Name
	}

	hashFunc := transform.FindHash(hashType)
	if hashFunc == nil {
		return false, fmt.Errorf("unexpected hash algorithm %s", hashType)
	}
	h, err := transform.ComputeHash(o.File, o.Path, hashFunc())
	if err != nil {
		return false, err
	}

	switch o.Operator {
	case tokenizer.Equals:
		result = h == o.Value
	case tokenizer.NotEquals:
		result = h != o.Value
	default:
		err = &ErrUnsupportedOperator{o.Attribute, o.Operator}
	}
	return result, err
}
