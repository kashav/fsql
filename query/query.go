package query

// Query represents an input query.
type Query struct {
	Attributes map[string]bool
	Sources    []string
	Conditions []Condition
}

// HasAttribute checks if the query's attribute map contains the provided
// attribute.
func (q *Query) HasAttribute(attribute string) bool {
	_, found := q.Attributes[attribute]
	return found
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute  string
	Comparator TokenType
	Value      string
}
