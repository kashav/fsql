package query

import (
	"errors"
	"fmt"
	"io"
	"os/user"
	"path/filepath"
	"strings"
)

// RunParser runs the parser on the input string and returns the parsed AST.
func RunParser(input string) (*Query, error) {
	return (&parser{}).parse(input)
}

var allAttributes = map[string]bool{
	"mode": true,
	"name": true,
	"size": true,
	"time": true,
}

type parser struct {
	tokenizer *Tokenizer
	current   *Token
	expected  TokenType
}

// Return true when no attributes are provided (regardless of if the SELECT
// keyword is provided). Returns false otherwise.
func (p *parser) showAllAttributes() (bool, error) {
	if p.expect(Select) == nil {
		if p.current == nil {
			return false, nil
		}

		if p.current.Type == From || p.current.Type == Where {
			return true, nil
		}

		if p.current.Type == Identifier {
			return false, nil
		}

		return false, p.currentError()
	}

	current := p.expect(Identifier)
	if current != nil {
		p.current = current
		return false, nil
	}

	return true, nil
}

// Parse each of the clauses in the input string.
func (p *parser) parse(input string) (*Query, error) {
	p.tokenizer = NewTokenizer(input)
	q := new(Query)
	q.Transformations = make(map[string][]Function)

	all, err := p.showAllAttributes()
	if err != nil {
		return nil, err
	}
	if all {
		q.Attributes = allAttributes
	} else {
		q.Attributes = make(map[string]bool)
		err := p.parseAttributes(&q.Attributes, &q.Transformations)
		if err != nil {
			return nil, err
		}
	}

	q.Sources = map[string][]string{
		"include": make([]string, 0),
		"exclude": make([]string, 0),
	}
	if p.expect(From) == nil {
		err := p.currentError()
		if p.expect(Identifier) != nil {
			return nil, err
		}
		q.Sources["include"] = append(q.Sources["include"], ".")
	} else {
		err := p.parseSources(&q.Sources)
		if err != nil {
			return nil, err
		}

		// Replace the tilde with the home directory in each source directory. This
		// is only required when the query is wrapped in quotes, since the shell
		// will automatically expand tildes otherwise.
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		for _, sourceType := range []string{"include", "exclude"} {
			for i, src := range q.Sources[sourceType] {
				if strings.Contains(src, "~") {
					q.Sources[sourceType][i] = filepath.Join(usr.HomeDir, src[1:])
				}
			}
		}
	}

	if p.expect(Where) == nil {
		err := p.currentError()
		if p.expect(Identifier) == nil {
			return q, nil
		}
		return nil, err
	}
	root, err := p.parseConditionTree()
	if err != nil {
		return nil, err
	}
	q.ConditionTree = root

	return q, nil
}

// Parse the list of attributes provided to the SELECT clause.
func (p *parser) parseAttributes(attributes *map[string]bool, transformations *map[string][]Function) error {
	attribute := p.expect(Identifier)
	if attribute == nil {
		return p.currentError()
	}
	if attribute.Raw == "*" || attribute.Raw == "all" {
		*attributes = allAttributes
	} else {
		p.current = attribute
		attribute, err := p.parseAttribute(transformations)
		if err != nil {
			return err
		}
		if _, ok := allAttributes[attribute.Raw]; !ok {
			return &ErrUnknownToken{attribute.Raw}
		}
		(*attributes)[attribute.Raw] = true

	}

	if p.expect(Comma) == nil {
		return nil
	}

	return p.parseAttributes(attributes, transformations)
}

// Parses all the transformation applied to given attribute recursively
func (p *parser) parseAttribute(transformations *map[string][]Function) (*Token, error) {
	identifier := p.expect(Identifier)
	var currFunction Function
	if identifier != nil {
		if p.expect(OpenParen) != nil {
			currFunction = Function{Name: identifier.Raw, Arguments: make([]string, 0)}
			attribute, err := p.parseAttribute(transformations)
			if attribute == nil {
				return attribute, err
			}
			for {
				if token := p.expect(Identifier); token != nil {
					currFunction.Arguments = append(currFunction.Arguments, token.Raw)
					continue
				} else if token := p.expect(Comma); token != nil {
					continue
				} else if token := p.expect(CloseParen); token != nil {
					if _, ok := (*transformations)[attribute.Raw]; !ok {
						(*transformations)[attribute.Raw] = make([]Function, 0)
					}
					(*transformations)[attribute.Raw] = append((*transformations)[attribute.Raw], currFunction)
					return attribute, err
				}
				return nil, p.currentError()
			}
		} else {
			if _, ok := allAttributes[identifier.Raw]; !ok {
				return nil, &ErrUnknownToken{identifier.Raw}
			}
			return identifier, nil
		}
	}
	return nil, p.currentError()
}

// Parse the list of directories passed to the FROM clause. Expects that
// the sources input has an "include" and "exclude" key.
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

// Parse the condition passed to the WHERE clause.
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

	if s.len() == 0 {
		return nil, p.currentError()
	}

	if s.len() > 1 {
		return nil, errors.New("failed to parse condition tree")
	}

	return s.pop(), nil
}

// Parse a single condition, made up of the negation, identifier (attribute),
// comparator, and value.
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

	return &Condition{
		Attribute:  attr.Raw,
		Comparator: comp,
		Value:      value.Raw,
		Negate:     negate,
	}, nil
}

// Returns the next token if it matches the expectation, nil otherwise.
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

// Returns the current error, based on the parser's current Token and the
// previously expected TokenType (set in expect).
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
	return fmt.Sprintf("Expected %s; got %s", e.Expected.String(), e.Actual.String())
}

// ErrUnknownToken represents an unknown token error.
type ErrUnknownToken struct {
	Raw string
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("Unknown token: %s", e.Raw)
}
