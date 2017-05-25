package parser

import (
	"errors"
	"os"

	"gopkg.in/oleiade/lane.v1"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

// parseConditionTree parses the condition tree passed to the WHERE clause.
func (p *parser) parseConditionTree() (*query.ConditionNode, error) {
	stack := lane.NewStack()
	errFailedToParse := errors.New("Failed to parse conditions")

	for {
		if p.current = p.tokenizer.Next(); p.current == nil {
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
				if p.parseSubquery(condition) != nil {
					return nil, errFailedToParse
				}
			}

			leaf := query.ConditionNode{Condition: condition}
			// If type assert fails, we assume the previous node was nil (i.e. not
			// a ConditionNode).
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

// parseNextCond parses the next condition of the query.
func (p *parser) parseNextCond() (*query.Condition, error) {
	cond := &query.Condition{}

	// Parse the NOT keyword.
	if p.expect(tokenizer.Not) != nil {
		cond.Negate = true
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
	cond.Attribute = attr.Raw
	cond.AttributeModifiers = modifiers

	// If this condition has modifiers, then p.current was unset while parsing
	// the modifier, se we set the current token manually.
	if len(modifiers) > 0 {
		p.current = p.tokenizer.Next()
	}
	if p.current == nil {
		return nil, p.currentError()
	}
	cond.Comparator = p.current.Type
	p.current = nil

	// Parse subquery of format `(...)`.
	if p.expect(tokenizer.OpenParen) != nil {
		token := p.expect(tokenizer.Subquery)
		if token == nil {
			return nil, p.currentError()
		}
		cond.IsSubquery = true
		cond.Value = token.Raw
		if p.expect(tokenizer.CloseParen) == nil {
			return nil, p.currentError()
		}
		return cond, nil
	}

	// Parse list of values of format `[...]`.
	if p.expect(tokenizer.OpenBracket) != nil {
		values := make([]string, 0)
		for {
			if token := p.expect(tokenizer.Identifier); token != nil {
				values = append(values, token.Raw)
			}
			if p.expect(tokenizer.Comma) != nil {
				continue
			}
			if p.expect(tokenizer.CloseBracket) != nil {
				break
			}
		}
		cond.Value = values
		return cond, nil
	}

	// Not a list nor a subquery -> plain identifier!
	if token := p.expect(tokenizer.Identifier); token != nil {
		cond.Value = token.Raw
	} else {
		return nil, p.currentError()
	}

	return cond, nil
}

// parseSubquery parses a subquery by recursively evaluating it's condition(s).
// If the subquery contains references to aliases from the superquery, it's
// Subquery attribute is set. Otherwise, we evaluate it's Subquery and set
// it's Value to the result.
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
