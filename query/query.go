package query

import (
	"os"
	"path/filepath"
)

// Query represents an input query.
type Query struct {
	Attributes map[string]bool
	Modifiers  map[string][]Modifier

	Sources       map[string][]string
	SourceAliases map[string]string

	ConditionTree *ConditionNode
}

// NewQuery returns a pointer to a Query.
func NewQuery() *Query {
	return &Query{
		Attributes: make(map[string]bool, 0),
		Modifiers:  make(map[string][]Modifier, 0),
		Sources: map[string][]string{
			"include": make([]string, 0),
			"exclude": make([]string, 0),
		},
		SourceAliases: make(map[string]string, 0),
		ConditionTree: nil,
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
func (q *Query) Execute(workFunc interface{}) error {
	seen := make(map[string]bool)
	excluder := &RegexpExclude{exclusions: q.Sources["exclude"]}

	for _, src := range q.Sources["include"] {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == "." {
				return nil
			}

			// Avoid walking a single directory more than once.
			if _, ok := seen[path]; ok {
				return nil
			}
			seen[path] = true

			if excluder.ShouldExclude(path) {
				return nil
			}

			if !q.ConditionTree.evaluateTree(path, info) {
				return nil
			}

			results := q.applyModifiers(path, info)
			workFunc.(func(string, os.FileInfo, map[string]interface{}))(path, info, results)
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
