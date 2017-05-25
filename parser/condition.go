package parser

import (
	"errors"
	"os"

	"gopkg.in/oleiade/lane.v1"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

// ParseConditionTree parses the condition tree passed to the WHERE clause.
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
			condition, err := p.parseNextCond()
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

// ParseNextCond parses the next condition of the query.
func (p *parser) parseNextCond() (*query.Condition, error) {
	negate := false
	if p.expect(tokenizer.Not) != nil {
		negate = true
	}

	ident := p.expect(tokenizer.Identifier)
	if ident == nil {
		return nil, p.currentError()
	}
	p.current = ident

	var modifiers []query.Modifier
	attr, err := p.parseAttrModifiers(&modifiers)
	if err != nil {
		return nil, err
	}
	if attr == nil {
		return nil, p.currentError()
	}

	// If this condition has modifiers, then p.current has been unset in the
	// modifier parsing process, se we get the next token manually.
	if len(modifiers) > 0 {
		p.current = p.tokenizer.Next()
	}
	if p.current == nil {
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
		Attribute:          attr.Raw,
		AttributeModifiers: modifiers,
		Comparator:         comp,
		Value:              value.Raw,
		Negate:             negate,
		IsSubquery:         subquery,
		Subquery:           nil,
	}, nil
}

// ParseSubquery parses a subquery by recursively evaluating it's condition(s).
// If the subquery contains references to aliases from the superquery, it's
// Subquery attribute is set. Otherwise, it's subquery is evaluated, it's Value
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

	value := make(map[interface{}]bool, 0)

	q.Execute(func(path string, info os.FileInfo,
		results map[string]interface{}) {
		if q.HasAttribute("name") {
			value[results["name"]] = true
		} else if q.HasAttribute("size") {
			value[results["size"]] = true
		} else if q.HasAttribute("time") {
			value[results["time"]] = true
		} else if q.HasAttribute("mode") {
			value[results["mode"]] = true
		}
	})

	condition.Value = value
	condition.IsSubquery = false
	return nil
}
