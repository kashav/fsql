package query

import (
	"fmt"
	"os"
	"bytes"
)

// Function entity used for defining transformations on attributes
type Function struct {
	Name string
	Arguments []string
}


// Query represents an input query.
type Query struct {
	Attributes    map[string]bool
	Sources       map[string][]string
	ConditionTree *ConditionNode // Root node of this query's condition tree.
	Transformations map[string][]Function
}

// PrintTransformations function prints all transformations parsed from given query
func (q *Query) PrintTransformations() {
	for _, functions := range q.Transformations{
		for _, function := range functions{
			fmt.Println(function.String())
		}
	}
}

// Return string representation of a given function
func (f *Function) String() string{
	var buffer bytes.Buffer
	buffer.WriteString(f.Name)
	for _,v := range f.Arguments{
		buffer.WriteString(fmt.Sprintf(", %s",v))
	}
	return buffer.String()
}

// HasAttribute checks if the query's attribute map contains the provided
// attribute.
func (q *Query) HasAttribute(attributes ...string) bool {
	for _, attribute := range attributes {
		if _, found := q.Attributes[attribute]; found {
			return true
		}
	}
	return false
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
