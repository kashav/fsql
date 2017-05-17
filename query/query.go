package query

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Query represents an input query.
type Query struct {
	Attributes    map[string]bool
	Sources       map[string][]string
	ConditionTree *ConditionNode // Root node of this query's condition tree.
}

// ReduceInclusions reduces this query's sources by removing any source
// which is a subdirectory of another source.
func (q *Query) ReduceInclusions() error {
	redundants := make(map[int]bool, len(q.Sources["include"])-1)

	for i, base := range q.Sources["include"] {
		for j, target := range q.Sources["include"][i+1:] {
			if i == (i + j + 1) {
				break
			}

			if base == target {
				// Duplicate source entry.
				redundants[i+j+1] = true
				continue
			}

			rel, err := filepath.Rel(base, target)
			if err != nil || (rel[:2] == ".." && rel[len(rel)-1] != '.') {
				// filepath.Rel only returns error when can't make target relative to
				// base, i.e. they're disjoint (which is what we want).
				continue
			} else if strings.Contains(rel, "..") {
				// Base directory is redundant, we can exit the inner loop.
				redundants[i] = true
				break
			} else {
				// Target directory is redundant.
				redundants[i+j+1] = true
			}
		}
	}

	sources := make([]string, 0)
	for i := 0; i < len(q.Sources["include"]); i++ {
		// Skip all redundant directories.
		if _, ok := redundants[i]; ok {
			continue
		}

		// Return error iff directory doesn't exist. Should we just ignore
		// nonexistent directories instead?
		path := q.Sources["include"][i]
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("no such file or directory: %s", path)
		}
		sources = append(sources, q.Sources["include"][i])
	}
	q.Sources["include"] = sources
	return nil
}

// HasAttribute checks if the query's attribute map contains the provided
// attribute.
func (q *Query) HasAttribute(attribute string) bool {
	_, found := q.Attributes[attribute]
	return found
}

// ConditionNode represents a single node of a query's WHERE clause tree.
type ConditionNode struct {
	Type      TokenType
	Condition *Condition
	Left      *ConditionNode
	Right     *ConditionNode
}

func (root *ConditionNode) String() string {
	if root == nil {
		return "(nil)"
	}

	if root.Condition != nil {
		return fmt.Sprintf("(%s)", (*root).Condition.String())
	}

	return fmt.Sprintf("(%s (%s, %s))", root.Type, root.Left.String(),
		root.Right.String())
}

// Evaluate runs pre-order traversal on the ConditionNode tree rooted at root
// and evaluates each conditional along the path.
func (root *ConditionNode) Evaluate(file os.FileInfo, compareFn interface{}) bool {
	if root == nil {
		return true
	}

	if root.Condition != nil {
		return compareFn.(func(Condition, os.FileInfo) bool)(*root.Condition, file)
	}

	if root.Type == And {
		return root.Left.Evaluate(file, compareFn) &&
			root.Right.Evaluate(file, compareFn)
	}

	if root.Type == Or {
		if root.Left.Evaluate(file, compareFn) {
			return true
		}
		return root.Right.Evaluate(file, compareFn)
	}

	return false
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute  string
	Comparator TokenType
	Value      string
	Negate     bool
}

func (c *Condition) String() string {
	return fmt.Sprintf(
		"{attribute: %s, comparator: %s, value: \"%s\", negate: %t}",
		c.Attribute, c.Comparator, c.Value, c.Negate)
}
