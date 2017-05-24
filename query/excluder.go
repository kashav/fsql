package query

import (
	"fmt"
	"regexp"
	"strings"
)

//Excluder allows us to support different methods of excluding in the future
type Excluder interface {
	ShouldExclude(path string) bool
}

//RegexpExclude uses regular expressions to tell if a file/path should be excluded
type RegexpExclude struct {
	recursive  bool
	exclusions []string
	regex      *regexp.Regexp
}

//ShouldExclude will return a boolean denoting whether or not the path should be excluded based on the given slice of exclusions
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
		if strings.HasSuffix(exclusion, "/") {
			tmpExclusions[i] = fmt.Sprintf(
				strings.Replace("^"+strings.TrimRight(exclusion, "/"), ".", "\\.", -1)) + "(/.*)?$"
		} else {
			tmpExclusions[i] = fmt.Sprintf(
				strings.Replace("^"+exclusion, ".", "\\.", -1)) + "(/.*)?$"
		}
	}
	fmt.Println(strings.Join(tmpExclusions, "|"))
	r.regex = regexp.MustCompile(strings.Join(tmpExclusions, "|"))
}
