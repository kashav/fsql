package query

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

// cmpAlpha performs alphabetic comparison on a and b.
func cmpAlpha(comp tokenizer.TokenType, a, b interface{}) bool {
	switch comp {
	case tokenizer.Equals:
		return a.(string) == b.(string)
	case tokenizer.NotEquals:
		return a.(string) != b.(string)
	case tokenizer.Like:
		aStr, bStr := a.(string), b.(string)
		if bStr[0] == '%' && bStr[len(bStr)-1] == '%' {
			return strings.Contains(aStr, bStr[1:len(bStr)-1])
		}
		if bStr[0] == '%' {
			return strings.HasSuffix(aStr, bStr[1:])
		}
		if bStr[len(bStr)-1] == '%' {
			return strings.HasPrefix(aStr, bStr[:len(bStr)-1])
		}
		return strings.Contains(aStr, bStr)
	case tokenizer.RLike:
		return regexp.MustCompile(b.(string)).MatchString(a.(string))
	case tokenizer.In:
		switch t := b.(type) {
		case map[interface{}]bool:
			if _, ok := t[a.(string)]; ok {
				return true
			}
		case []string:
			for _, el := range t {
				if a.(string) == el {
					return true
				}
			}
		case string:
			for _, el := range strings.Split(t, ",") {
				if a.(string) == el {
					return true
				}
			}
		}
	}
	return false
}

// cmpNumeric performs numeric comparison on a and b.
func cmpNumeric(comp tokenizer.TokenType, a, b interface{}) bool {
	switch comp {
	case tokenizer.Equals:
		return a.(int64) == b.(int64)
	case tokenizer.NotEquals:
		return a.(int64) != b.(int64)
	case tokenizer.GreaterThanEquals:
		return a.(int64) >= b.(int64)
	case tokenizer.GreaterThan:
		return a.(int64) > b.(int64)
	case tokenizer.LessThanEquals:
		return a.(int64) <= b.(int64)
	case tokenizer.LessThan:
		return a.(int64) < b.(int64)
	case tokenizer.In:
		if _, ok := b.(map[interface{}]bool)[a.(int64)]; ok {
			return true
		}
	}
	return false
}

// cmpTime performs time comparison on a and b.
func cmpTime(comp tokenizer.TokenType, a, b interface{}) bool {
	switch comp {
	case tokenizer.Equals:
		return a.(time.Time).Equal(b.(time.Time))
	case tokenizer.NotEquals:
		return !a.(time.Time).Equal(b.(time.Time))
	case tokenizer.GreaterThanEquals:
		return a.(time.Time).After(b.(time.Time)) || a.(time.Time).Equal(b.(time.Time))
	case tokenizer.GreaterThan:
		return a.(time.Time).After(b.(time.Time))
	case tokenizer.LessThanEquals:
		return a.(time.Time).Before(b.(time.Time)) || a.(time.Time).Equal(b.(time.Time))
	case tokenizer.LessThan:
		return a.(time.Time).Before(b.(time.Time))
	case tokenizer.In:
		if _, ok := b.(map[interface{}]bool)[a.(time.Time)]; ok {
			return true
		}
	}
	return false
}

// cmpMode performs mode comparison with file and fileType.
func cmpMode(comp tokenizer.TokenType, file os.FileInfo, fileType interface{}) bool {
	if comp != tokenizer.Is {
		return false
	}
	switch strings.ToUpper(fileType.(string)) {
	case "DIR":
		return file.Mode().IsDir()
	case "REG":
		return file.Mode().IsRegular()
	}
	return false
}
