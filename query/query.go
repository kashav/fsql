package query

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	uBYTE     = 1.0
	uKILOBYTE = 1024 * uBYTE
	uMEGABYTE = 1024 * uKILOBYTE
	uGIGABYTE = 1024 * uMEGABYTE
)

// Query represents an input query.
type Query struct {
	Attributes    map[string]bool
	Sources       map[string][]string
	ConditionTree *ConditionNode
	SourceAliases map[string]string
}

// NewQuery returns a pointer to a Query.
func NewQuery() *Query {
	return &Query{
		Attributes: make(map[string]bool, 0),
		Sources: map[string][]string{
			"include": make([]string, 0),
			"exclude": make([]string, 0),
		},
		ConditionTree: nil,
		SourceAliases: make(map[string]string, 0),
	}
}

// HasAttribute checks if this query contains any of the provided attributes.
func (q *Query) HasAttribute(attributes ...string) bool {
	for _, attribute := range attributes {
		if _, found := q.Attributes[attribute]; found {
			return true
		}
	}
	return false
}

// Execute runs the query by walking the full path of each source and
// evaluating the condition tree for each file. This method calls workFunc on
// each "successful" file.
func (q *Query) Execute(workFunc interface{}) {
	containsAny := func(path string) bool {
		for _, exclusion := range q.Sources["exclude"] {
			if strings.Contains(path, exclusion) {
				return true
			}
		}

		return false
	}

	seen := make(map[string]bool)

	for _, src := range q.Sources["include"] {
		filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if path == "." || path == ".." || err != nil {
				return nil
			}

			// Avoid walking a single directory more than once.
			if _, ok := seen[path]; ok {
				return nil
			}
			seen[path] = true

			if containsAny(path) || !q.ConditionTree.evaluateTree(path, info) {
				return nil
			}

			workFunc.(func(string, os.FileInfo))(path, info)
			return nil
		})
	}
}
