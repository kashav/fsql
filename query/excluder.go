package query

import (
	"regexp"
	"strings"
)

type Excluder interface {
	ShouldExclude(path string) bool
}

//RegexpExclude is a struct that will decide if exclusions should be excluded
type RegexpExclude struct {
	recursive  bool
	exclusions []string
}

//ShouldExclude will return a boolean denoting whether or not the path should be excluded based on the user input
func (r RegexpExclude) ShouldExclude(path string) bool {
	regex := r.buildRegex()
	if b, ok := regexp.MatchString(regex, path); ok == nil {
		return b
	}
	return false
}

func (r RegexpExclude) buildRegex() string {
	var regex string
	var prev string
	var curr string
	for _, exclusion := range r.exclusions {
		prev = curr
		if strings.HasSuffix(exclusion, "/") {
			curr =
				or(mustEnd(mustBegin(escape(exclusion))),
					mustEnd(mustBegin(escape(exclusion[:len(exclusion)-1]))))
		} else {
			curr =
				or(mustEnd(mustBegin(escape(exclusion))),
					mustEnd(mustBegin(escape(exclusion+"/"))))
		}
		regex = or(prev, curr)
	}
	return regex
}

func or(p1, p2 string) string {
	if p1 == "" {
		return p2
	}
	if p2 == "" {
		return p1
	}
	return p1 + "|" + p2
}

func mustEnd(p string) string {
	return p + "$"
}

func mustBegin(p string) string {
	return "^" + p
}

func group(expression string) string {
	return "(" + expression + ")"
}

func escape(expression string) string {
	var str string
	for _, r := range expression {
		if r == '.' {
			str += "\\"
		}
		str += string(r)
	}

	return str
}
