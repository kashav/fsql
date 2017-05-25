package parser

import (
	"fmt"

	"github.com/kshvmdn/fsql/tokenizer"
)

// ErrUnexpectedToken represents an unexpected token error.
type ErrUnexpectedToken struct {
	Actual   tokenizer.TokenType
	Expected tokenizer.TokenType
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("expected %s; got %s", e.Expected.String(),
		e.Actual.String())
}

// ErrUnknownToken represents an unknown token error.
type ErrUnknownToken struct {
	Raw string
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("unknown token: %s", e.Raw)
}
