package query

import (
	"os"
	"path/filepath"
	"strings"
)

// Query represents an input query.
type Query struct {
	Attributes []string
	Modifiers  map[string][]Modifier

	Sources       map[string][]string
	SourceAliases map[string]string

	ConditionTree *ConditionNode
}

// NewQuery returns a pointer to a Query.
func NewQuery() *Query {
	return &Query{
		Attributes: make([]string, 0),
		Modifiers:  make(map[string][]Modifier),
		Sources: map[string][]string{
			"include": make([]string, 0),
			"exclude": make([]string, 0),
		},
		SourceAliases: make(map[string]string),
		ConditionTree: nil,
	}
}

// HasAttribute checks if this query contains any of the provided attributes.
func (q *Query) HasAttribute(attributes ...string) bool {
	for _, attribute := range attributes {
		for _, queryAttribute := range q.Attributes {
			if attribute == queryAttribute {
				return true
			}
		}
	}
	return false
}

// Execute runs the query by walking the full path of each source and
// evaluating the condition tree for each file. This method calls workFunc on
// each "successful" file.
func (q *Query) Execute(workFunc interface{}) error {
	seen := map[string]bool{}
	excluder := &regexpExclude{exclusions: q.Sources["exclude"]}

	for _, src := range q.Sources["include"] {
		// TODO: Improve our method of detecting if src is a glob pattern. This
		// currently doesn't support usage of square brackets, since the tokenizer
		// doesn't recognize these as part of a directory.
		//
		// Pattern reference: https://golang.org/pkg/path/filepath/#Match.
		if strings.ContainsAny(src, "*?") {
			// If src does _resemble_ a glob pattern, we find all matches and
			// evaluate the condition tree against each.
			matches, err := filepath.Glob(src)
			if err != nil {
				return err
			}

			for _, match := range matches {
				if err = filepath.Walk(match, q.walkFunc(seen, excluder, workFunc)); err != nil {
					return err
				}
			}
			continue
		}

		if err := filepath.Walk(src, q.walkFunc(seen, excluder, workFunc)); err != nil {
			return err
		}
	}

	return nil
}

// walkFunc returns a filepath.WalkFunc which evaluates the condition tree
// against the given file.
func (q *Query) walkFunc(seen map[string]bool, excluder Excluder,
	workFunc interface{}) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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

		if excluder.shouldExclude(path) {
			return nil
		}

		if ok, err := q.ConditionTree.evaluateTree(path, info); err != nil {
			return err
		} else if !ok {
			return nil
		}

		results, err := q.applyModifiers(path, info)
		if err != nil {
			return err
		}
		workFunc.(func(string, os.FileInfo, map[string]interface{}))(path, info, results)
		return nil
	}
}
