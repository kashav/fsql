package evaluate

import (
	"fmt"

	"github.com/kashav/fsql/tokenizer"
)

// ErrUnsupportedAttribute represents an unsupported attribute error.
type ErrUnsupportedAttribute struct {
	Attribute string
}

func (e *ErrUnsupportedAttribute) Error() string {
	return fmt.Sprintf("unsupported attribute %s", e.Attribute)
}

// ErrUnsupportedType represents an unsupported type error.
type ErrUnsupportedType struct {
	Attribute string
	Value     interface{}
}

func (e *ErrUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type %T for for attribute %s", e.Value,
		e.Attribute)
}

// ErrUnsupportedOperator represents an unsupported operator error.
type ErrUnsupportedOperator struct {
	Attribute string
	Operator  tokenizer.TokenType
}

func (e *ErrUnsupportedOperator) Error() string {
	return fmt.Sprintf("unsupported operator %s for attribute %s",
		e.Operator.String(), e.Attribute)
}
