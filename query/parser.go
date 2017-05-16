package query

import (
	"errors"
	"fmt"
	"io"
	"strings"
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
	q.Sources = map[string][]string{
		"include": make([]string, 0),
		"exclude": make([]string, 0),
	}
	err = p.parseSources(&q.Sources)
	if err != nil {
		return nil, err
	}

	if p.expect(Where) == nil {
		return nil, p.currentError()
	}

	root, err := p.parseConditionTree()
	if err != nil {
		return nil, err
	}
	q.ConditionTree = root

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

func (p *parser) parseSources(sources *map[string][]string) error {
	sourceType := "include"
	if p.expect(Minus) != nil {
		sourceType = "exclude"
	}
	source := p.expect(Identifier)
	if source == nil {
		return p.currentError()
	}
	(*sources)[sourceType] = append((*sources)[sourceType], source.Raw)
	if p.expect(Comma) == nil {
		return nil
	}
	return p.parseSources(sources)
}

func (p *parser) parseConditionTree() (*ConditionNode, error) {
	s := new(stack)

	for {
		p.current = p.tokenizer.Next()
		if p.current == nil {
			break
		}

		switch p.current.Type {
		case Not:
			fallthrough
		case Identifier:
			condition, err := p.parseNextCondition()
			if err != nil {
				return nil, p.currentError()
			}

			leaf := ConditionNode{Condition: condition}
			previous := s.pop()
			if previous == nil {
				s.push(&leaf)
			} else {
				if (*previous).Condition == nil {
					(*previous).Right = &leaf
				}
				s.push(previous)
			}
		case And:
			fallthrough
		case Or:
			left := s.pop()
			node := ConditionNode{
				Type: p.current.Type,
				Left: left,
			}
			s.push(&node)
		case OpenParen:
			s.push(nil)
		case CloseParen:
			right := s.pop()
			root := s.pop()
			if root != nil {
				root.Right = right
				s.push(root)
			} else {
				s.push(right)
			}
		}
	}

	if s.len() > 1 {
		return nil, errors.New("failed to parse condition tree")
	}

	return s.pop(), nil
}

func (p *parser) parseNextCondition() (*Condition, error) {
	negate := false
	if p.expect(Not) != nil {
		negate = true
	}

	attr := p.expect(Identifier)
	if attr == nil {
		return nil, p.currentError()
	}

	p.current = p.tokenizer.Next()
	if p.current == nil {
		return nil, p.currentError()
	}
	comp := p.current.Type
	p.current = nil

	value := p.expect(Identifier)
	if value == nil {
		return nil, p.currentError()
	}

	// Use regexp to evaluate wildcard (%) in LIKE conditions.
	if comp == Like && strings.Contains(value.Raw, "%") {
		comp = RLike
		value.Raw = strings.Replace(value.Raw, "%", ".*", -1)
	}

	return &Condition{
		Attribute:  attr.Raw,
		Comparator: comp,
		Value:      value.Raw,
		Negate:     negate,
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
