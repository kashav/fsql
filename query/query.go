package query

import (
	"fmt"
	"os"
)

// Query represents an input query.
type Query struct {
	Attributes    map[string]bool
	Sources       map[string][]string
	ConditionTree *ConditionNode // Root node of this query's condition tree.
}

// ConditionNode represents a single node of a query's WHERE clause tree.
type ConditionNode struct {
	Type      TokenType
	Condition *Condition
	Left      *ConditionNode
	Right     *ConditionNode
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute  string
	Comparator TokenType
	Value      string
	Negate     bool
}

// HasAttribute checks if the query's attribute map contains the provided
// attribute.
func (q *Query) HasAttribute(attribute string) bool {
	_, found := q.Attributes[attribute]
	return found
}

func (root *ConditionNode) String() string {
	if root == nil {
		return "<nil>"
	}

	if root.Condition != nil {
		return fmt.Sprintf("<%s>", *root.Condition)
	}

	return fmt.Sprintf("<%s (%s, %s)>", root.Type, root.Left, root.Right)
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

	left := root.Left.Evaluate(file, compareFn)
	right := root.Right.Evaluate(file, compareFn)

	if root.Type == And {
		return left && right
	}

	if root.Type == Or {
		return left || right
	}

	return false
}
