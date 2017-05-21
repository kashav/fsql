package query

import (
	"fmt"
	"regexp"
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
func (r *RegexpExclude) ShouldExclude(path string) bool {
	regex := r.buildRegex()
	if b, ok := regexp.MatchString(regex, path); ok == nil {
		return b
	}
	return false
}

func (r *RegexpExclude) buildRegex() string {
	for _, exclusion := range r.exclusions {
		fmt.Println(exclusion)
	}
	return ""
}

func or(p1, p2 string) string {
	return p1 + "|" + p2
}

func matchEnd(p string) string {
	return p + "$"
}

func matchBegin(p string) string {
	return "^" + p
}

func group(expression string) string {
	return "(" + expression + ")"

}
