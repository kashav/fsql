package parser

import (
	"fmt"
	"io"

	"github.com/kashav/fsql/tokenizer"
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

// currentError returns the current error, based on the parser's current Token
// and the previously expected TokenType (set in parser.expect).
func (p *parser) currentError() error {
	if p.current == nil {
		return io.ErrUnexpectedEOF
	}

	if p.current.Type == tokenizer.Unknown {
		return &ErrUnknownToken{Raw: p.current.Raw}
	}

	return &ErrUnexpectedToken{Actual: p.current.Type, Expected: p.expected}
}
