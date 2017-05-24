package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/oleiade/lane.v1"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

var allAttributes = map[string]bool{
	"mode": true,
	"name": true,
	"size": true,
	"time": true,
}

// Run parses the input string and returns the parsed AST (query).
func Run(input string) (*query.Query, error) {
	return (&parser{}).parse(input)
}

type parser struct {
	tokenizer *tokenizer.Tokenizer
	current   *tokenizer.Token
	expected  tokenizer.TokenType
}

// Parse each of the clauses in the input string.
func (p *parser) parse(input string) (*query.Query, error) {
	p.tokenizer = tokenizer.NewTokenizer(input)
	q := query.NewQuery()

	if err := p.parseSelectClause(q); err != nil {
		return nil, err
	}

	if err := p.parseFromClause(q); err != nil {
		return nil, err
	}

	if err := p.parseWhereClause(q); err != nil {
		return nil, err
	}

	return q, nil
}

// Parse the SELECT clause.
func (p *parser) parseSelectClause(q *query.Query) error {
	// Determine if we should show all attributes. This is only true when
	// no attributes are provided (regardless of if the SELECT keyword is
	// provided or not).
	showAll := true
	if p.expect(tokenizer.Select) == nil {
		if p.current == nil || p.current.Type == tokenizer.Identifier {
			showAll = false
		} else if p.current.Type == tokenizer.From || p.current.Type == tokenizer.Where {
			// No SELECT and next token is FROM/WHERE, show all!
			showAll = true
		} else {
			// No SELECT and next token is not Ident nor FROM/WHERE -> malformed input.
			return p.currentError()
		}
	} else {
		if current := p.expect(tokenizer.Identifier); current != nil {
			p.current = current
			showAll = false
		}
	}

	if showAll {
		q.Attributes = allAttributes
	} else {
		err := p.parseAttrList(&q.Attributes, &q.Modifiers)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseAttrList parses the list of attributes passed to the SELECT clause.
func (p *parser) parseAttrList(attributes *map[string]bool,
	modifiers *map[string][]query.Modifier) error {
	attribute := p.expect(tokenizer.Identifier)
	if attribute == nil {
		return p.currentError()
	}

	if attribute.Raw == "*" || attribute.Raw == "all" {
		*attributes = allAttributes
	} else {
		p.current = attribute
		attribute, err := p.parseSingleAttr(modifiers)
		if err != nil {
			return err
		}

		if _, ok := allAttributes[attribute.Raw]; !ok {
			return &ErrUnknownToken{attribute.Raw}
		}

		(*attributes)[attribute.Raw] = true
	}

	// If next token is a comma, recurse!
	if p.expect(tokenizer.Comma) != nil {
		return p.parseAttrList(attributes, modifiers)
	}

	return nil
}

// ParseSingleAttr parses a single attribute and it's associated modifiers.
func (p *parser) parseSingleAttr(modifiers *map[string][]query.Modifier) (*tokenizer.Token, error) {
	var current query.Modifier

	ident := p.expect(tokenizer.Identifier)
	if ident == nil {
		return nil, p.currentError()
	}

	if p.expect(tokenizer.OpenParen) == nil {
		// No modifier on this attribute
		if _, ok := allAttributes[ident.Raw]; !ok {
			return nil, &ErrUnknownToken{ident.Raw}
		}
		return ident, nil
	}

	current = query.Modifier{
		Name:      ident.Raw,
		Arguments: make([]string, 0),
	}

	attribute, err := p.parseSingleAttr(modifiers)
	if attribute == nil || err != nil {
		return nil, err
	}

	for {
		if token := p.expect(tokenizer.Identifier); token != nil {
			current.Arguments = append(current.Arguments, token.Raw)
			continue
		}

		if token := p.expect(tokenizer.Comma); token != nil {
			continue
		}

		if token := p.expect(tokenizer.CloseParen); token != nil {
			if _, ok := (*modifiers)[attribute.Raw]; !ok {
				(*modifiers)[attribute.Raw] = make([]query.Modifier, 0)
			}
			(*modifiers)[attribute.Raw] = append((*modifiers)[attribute.Raw], current)
			return attribute, nil
		}
	}
}

// Parse the FROM clause.
func (p *parser) parseFromClause(q *query.Query) error {
	if p.expect(tokenizer.From) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) != nil {
			// No FROM, but an identifier -> malformed query.
			return err
		}
		q.Sources["include"] = append(q.Sources["include"], ".")
	} else {
		err := p.parseSources(&q.Sources, &q.SourceAliases)
		if err != nil {
			return err
		}

		// Replace the tilde with the home directory in each source directory. This
		// is only required when the query is wrapped in quotes, since the shell
		// will automatically expand tildes otherwise.
		usr, err := user.Current()
		if err != nil {
			return err
		}
		for _, sourceType := range []string{"include", "exclude"} {
			for i, src := range q.Sources[sourceType] {
				if strings.Contains(src, "~") {
					q.Sources[sourceType][i] = filepath.Join(usr.HomeDir, src[1:])
				}
			}
		}
	}

	return nil
}

// Parse the list of directories passed to the FROM clause. If a source is
// followed by the AS keyword, the following word is registered as an alias.
func (p *parser) parseSources(sources *map[string][]string, aliases *map[string]string) error {
	sourceType := "include"
	if p.expect(tokenizer.Hyphen) != nil {
		sourceType = "exclude"
	}

	source := p.expect(tokenizer.Identifier)
	if source == nil {
		return p.currentError()
	}
	(*sources)[sourceType] = append((*sources)[sourceType], source.Raw)

	if p.expect(tokenizer.As) != nil {
		alias := p.expect(tokenizer.Identifier)
		if alias == nil {
			return p.currentError()
		}
		(*aliases)[alias.Raw] = source.Raw
	}

	// If next token is a comma, recurse!
	if p.expect(tokenizer.Comma) != nil {
		return p.parseSources(sources, aliases)
	}

	return nil
}

// Parse the WHERE clause.
func (p *parser) parseWhereClause(q *query.Query) error {
	if p.expect(tokenizer.Where) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) == nil {
			return nil
		}
		return err
	}
	root, err := p.parseConditionTree()
	if err != nil {
		return err
	}
	q.ConditionTree = root

	return nil
}

// Parse the condition tree passed to the WHERE clause.
func (p *parser) parseConditionTree() (*query.ConditionNode, error) {
	stack := lane.NewStack()
	errFailedToParse := errors.New("Failed to parse conditions")

	for {
		p.current = p.tokenizer.Next()
		if p.current == nil {
			break
		}

		switch p.current.Type {
		case tokenizer.Not:
			fallthrough
		case tokenizer.Identifier:
			condition, err := p.parseNextCondition()
			if err != nil {
				return nil, p.currentError()
			}
			if condition.IsSubquery {
				if err := p.parseSubquery(condition); err != nil {
					return nil, errFailedToParse
				}
			}

			leaf := query.ConditionNode{Condition: condition}
			if prev, ok := stack.Pop().(*query.ConditionNode); !ok {
				stack.Push(&leaf)
			} else {
				if prev.Condition == nil {
					prev.Right = &leaf
				}
				stack.Push(prev)
			}
		case tokenizer.And:
			fallthrough
		case tokenizer.Or:
			left, ok := stack.Pop().(*query.ConditionNode)
			if !ok {
				return nil, errFailedToParse
			}

			node := query.ConditionNode{
				Type: p.current.Type,
				Left: left,
			}
			stack.Push(&node)
		case tokenizer.OpenParen:
			stack.Push(nil)
		case tokenizer.CloseParen:
			right, ok := stack.Pop().(*query.ConditionNode)
			if !ok {
				return nil, errFailedToParse
			}

			if root, ok := stack.Pop().(*query.ConditionNode); ok {
				root.Right = right
				stack.Push(root)
			} else {
				stack.Push(right)
			}
		}
	}

	if stack.Size() == 0 {
		return nil, p.currentError()
	}

	if stack.Size() > 1 {
		return nil, errFailedToParse
	}

	node, ok := stack.Pop().(*query.ConditionNode)
	if !ok {
		return nil, errFailedToParse
	}
	return node, nil
}

// Parse a single condition, made up of the negation, identifier (attribute),
// comparator, and value.
func (p *parser) parseNextCondition() (*query.Condition, error) {
	negate := false
	if p.expect(tokenizer.Not) != nil {
		negate = true
	}

	attr := p.expect(tokenizer.Identifier)
	if attr == nil {
		return nil, p.currentError()
	}

	p.current = p.tokenizer.Next()
	if p.current == nil {
		// FIXME: Return a more appropriate error.
		return nil, p.currentError()
	}
	comp := p.current.Type
	p.current = nil

	var value *tokenizer.Token
	var subquery bool
	if p.expect(tokenizer.OpenParen) != nil {
		value = p.expect(tokenizer.Subquery)
		subquery = true
	} else {
		value = p.expect(tokenizer.Identifier)
	}
	if value == nil {
		return nil, p.currentError()
	}

	// We check for a closing paren AFTER checking that value is non-nil to
	// prevent the current error from being overwritten.
	if subquery && p.expect(tokenizer.CloseParen) == nil {
		return nil, p.currentError()
	}

	return &query.Condition{
		Attribute:  attr.Raw,
		Comparator: comp,
		Value:      value.Raw,
		Negate:     negate,
		IsSubquery: subquery,
		Subquery:   nil,
	}, nil
}

// Parse a subquery by recursively evaluating it's condition(s). If the
// subquery contains references to aliases from the superquery, it's Subquery
// attribute is set. Otherwise, it's subquery is evaluated, it's Value
// attribute is set to the returned result list, and it's IsSubquery attribute
// is made false.
func (p *parser) parseSubquery(condition *query.Condition) error {
	q, err := Run(condition.Value.(string))
	if err != nil {
		return err
	}

	if len(q.SourceAliases) > 0 {
		condition.Subquery = q
		return nil
	}

	results := make(map[interface{}]bool, 0)

	q.Execute(func(path string, info os.FileInfo) {
		// FIXME: Allow more than 1 attribute? If so, how will we compare values?
		if q.HasAttribute("name") {
			results[info.Name()] = true
		} else if q.HasAttribute("size") {
			results[info.Size()] = true
		} else if q.HasAttribute("time") {
			results[info.ModTime()] = true
		} else if q.HasAttribute("mode") {
			results[info.Mode()] = true
		}
	})

	condition.Value = results
	condition.IsSubquery = false
	return nil
}

// Returns the next token if it matches the expectation, nil otherwise.
func (p *parser) expect(t tokenizer.TokenType) *tokenizer.Token {
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
// previously expected TokenType (set in parser.expect).
func (p *parser) currentError() error {
	if p.current == nil {
		return io.ErrUnexpectedEOF
	}

	if p.current.Type == tokenizer.Unknown {
		return &ErrUnknownToken{Raw: p.current.Raw}
	}

	return &ErrUnexpectedToken{Actual: p.current.Type, Expected: p.expected}
}

// ErrUnexpectedToken represents an unexpected token error.
type ErrUnexpectedToken struct {
	Actual   tokenizer.TokenType
	Expected tokenizer.TokenType
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("Expected: %s, got: %s",
		e.Expected.String(), e.Actual.String())
}

// ErrUnknownToken represents an unknown token error.
type ErrUnknownToken struct {
	Raw string
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("Unknown token: %s", e.Raw)
}
