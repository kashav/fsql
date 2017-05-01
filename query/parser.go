package query

import (
	"fmt"
	"io"
)

// Run parses input and returns the output.
func Run(input string) (*Query, error) {
	return (&parser{}).parse(input)
}

type parser struct {
	tokenizer *Tokenizer
	current   *Token
	expected  TokenType
}

func (p *parser) parse(input string) (*Query, error) {
	p.tokenizer = NewTokenizer(input)
	q := new(Query)

	if p.expect(Select) == nil {
		return nil, p.currentError()
	}
	q.Attributes = make(map[string]bool)
	err := p.parseAttributes(&q.Attributes)
	if err != nil {
		return nil, err
	}

	if p.expect(From) == nil {
		return nil, p.currentError()
	}
	sources, err := p.parseSources()
	if err != nil {
		return nil, err
	}
	q.Sources = sources

	if p.expect(Where) == nil {
		return nil, p.currentError()
	}
	conditions, err := p.parseConditions()
	if err != nil {
		return nil, err
	}
	q.Conditions = conditions

	return q, nil
}

func (p *parser) parseAttributes(attributes *map[string]bool) error {
	attribute := p.expect(Identifier)
	if attribute == nil {
		return p.currentError()
	}
	if attribute.Raw == "*" {
		*attributes = map[string]bool{
			"mode": true,
			"name": true,
			"size": true,
			"time": true,
		}
	} else {
		(*attributes)[attribute.Raw] = true
	}
	if p.expect(Comma) == nil {
		return nil
	}
	return p.parseAttributes(attributes)
}

func (p *parser) parseSources() ([]string, error) {
	sources := []string{}

	source := p.expect(Identifier)
	if source == nil {
		return nil, p.currentError()
	}
	sources = append(sources, source.Raw)

	for {
		if p.expect(Comma) == nil {
			return sources, nil
		}

		source = p.expect(Identifier)
		if source == nil {
			return nil, p.currentError()
		}
		sources = append(sources, source.Raw)
	}
}

func (p *parser) parseConditions() ([]Condition, error) {
	conditions := []Condition{}

	cond, err := p.parseNextCondition()
	if err != nil {
		return nil, err
	}
	conditions = append(conditions, *cond)

	for {
		if p.expect(Comma) == nil {
			return conditions, nil
		}

		cond, err = p.parseNextCondition()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, *cond)
	}
}

func (p *parser) parseNextCondition() (*Condition, error) {
	attr := p.expect(Identifier)
	if attr == nil {
		return nil, p.currentError()
	}

	p.current = p.tokenizer.Next()
	if p.current == nil {
		return nil, p.currentError()
	}
	// TODO: check that p.current is a valid comparator
	comp := p.current.Type
	p.current = nil

	value := p.expect(Identifier)
	if value == nil {
		return nil, p.currentError()
	}

	return &Condition{
		Attribute:  attr.Raw,
		Comparator: comp,
		Value:      value.Raw,
	}, nil
}

func (p *parser) expect(t TokenType) *Token {
	p.expected = t

	if p.current == nil {
		p.current = p.tokenizer.Next()
	}

	if p.current != nil && p.current.Type == t {
		tok := p.current
		p.current = nil
		return tok
	}

	return nil
}

func (p *parser) currentError() error {
	if p.current == nil {
		return io.ErrUnexpectedEOF
	}

	if p.current.Type == Unknown {
		return &ErrUnknownToken{Raw: p.current.Raw}
	}

	return &ErrUnexpectedToken{Actual: p.current.Type, Expected: p.expected}
}

// ErrUnexpectedToken represents an unexpected token error.
type ErrUnexpectedToken struct {
	Actual   TokenType
	Expected TokenType
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("expected %s; got %s", e.Expected.String(), e.Actual.String())
}

// ErrUnknownToken represents an unknown token error.
type ErrUnknownToken struct {
	Raw string
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("unknown token: %s", e.Raw)
}
