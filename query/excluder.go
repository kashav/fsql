package query

import (
	"fmt"
	"regexp"
	"strings"
)

// Excluder allows us to support different methods of excluding in the future.
type Excluder interface {
	ShouldExclude(path string) bool
}

// RegexpExclude uses regular expressions to tell if a file/path should be
// excluded.
type RegexpExclude struct {
	exclusions []string
	regex      *regexp.Regexp
}

// ShouldExclude will return a boolean denoting whether or not the path should
// be excluded based on the given slice of exclusions.
func (r *RegexpExclude) ShouldExclude(path string) bool {
	if r.regex == nil {
		r.buildRegex()
	}
	if r.regex.String() == "" {
		return false
	}
	return r.regex.MatchString(path)
}

func (r *RegexpExclude) buildRegex() {
	numExclusion := len(r.exclusions)
	tmpExclusions := make([]string, numExclusion, numExclusion)
	for i, exclusion := range r.exclusions {
		// Wrap exclusion in ^ and (/.*)?$ AFTER replacing trailing forward
		// slash and replacing all dots with `\\.`
		tmpExclusions[i] = fmt.Sprintf("^%s(/.*)?$",
			strings.Replace(strings.TrimRight(exclusion, "/"), ".", "\\.", -1))
	}
	r.regex = regexp.MustCompile(strings.Join(tmpExclusions, "|"))
}
