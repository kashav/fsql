package query

// Query represents an input query.
type Query struct {
	Attributes []string
	Sources    []string
	Conditions []Condition
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute  string
	Comparator TokenType
	Value      string
}
