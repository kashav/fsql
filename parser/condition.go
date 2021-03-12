package parser

import (
	"errors"
	"os"

	"github.com/oleiade/lane"

	"github.com/kashav/fsql/query"
	"github.com/kashav/fsql/tokenizer"
)

// parseConditionTree parses the condition tree passed to the WHERE clause.
func (p *parser) parseConditionTree() (*query.ConditionNode, error) {
	stack := lane.NewStack()
	errFailedToParse := errors.New("failed to parse conditions")

	for {
		if p.current = p.tokenizer.Next(); p.current == nil {
			break
		}

		switch p.current.Type {

		case tokenizer.Not:
			// TODO: Handle NOT (...), for the time being we proceed with the other
			// tokens and handle the negation when parsing the condition.
			fallthrough

		case tokenizer.Identifier:
			condition, err := p.parseCondition()
			if err != nil {
				return nil, err
			}
			if condition == nil {
				return nil, p.currentError()
			}

			if condition.IsSubquery {
				if err := p.parseSubquery(condition); err != nil {
					return nil, err
				}
			}

			leafNode := &query.ConditionNode{Condition: condition}
			if prevNode, ok := stack.Pop().(*query.ConditionNode); !ok {
				stack.Push(leafNode)
			} else if prevNode.Condition == nil {
				prevNode.Right = leafNode
				stack.Push(prevNode)
			} else {
				return nil, errFailedToParse
			}

		case tokenizer.And, tokenizer.Or:
			leftNode, ok := stack.Pop().(*query.ConditionNode)
			if !ok {
				return nil, errFailedToParse
			}

			node := query.ConditionNode{
				Type: &p.current.Type,
				Left: leftNode,
			}
			stack.Push(&node)

		case tokenizer.OpenParen:
			stack.Push(nil)

		case tokenizer.CloseParen:
			rightNode, ok := stack.Pop().(*query.ConditionNode)
			if !ok {
				return nil, errFailedToParse
			}

			if rootNode, ok := stack.Pop().(*query.ConditionNode); !ok {
				stack.Push(rightNode)
			} else {
				rootNode.Right = rightNode
				stack.Push(rootNode)
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

// parseCondition parses and returns the next condition.
func (p *parser) parseCondition() (*query.Condition, error) {
	cond := &query.Condition{}

	// If we find a NOT, negate the condition.
	if p.expect(tokenizer.Not) != nil {
		cond.Negate = true
	}

	ident := p.expect(tokenizer.Identifier)
	if ident == nil {
		return nil, p.currentError()
	}
	p.current = ident

	var modifiers []query.Modifier
	attr, err := p.parseAttr(&modifiers)
	if err != nil {
		return nil, err
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
	cond.Operator = p.current.Type
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
	token := p.expect(tokenizer.Identifier)
	if token == nil {
		return nil, p.currentError()
	}
	cond.Value = token.Raw
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

	// If the subquery has aliases, we'll have to parse the subquery against
	// each file, so we don't do anything here.
	if len(q.SourceAliases) > 0 {
		condition.Subquery = q
		return nil
	}

	value := make(map[interface{}]bool, 0)
	workFunc := func(path string, info os.FileInfo, res map[string]interface{}) {
		for _, attr := range [...]string{"name", "size", "time", "mode"} {
			if q.HasAttribute(attr) {
				value[res[attr]] = true
				return
			}
		}
	}

	if err = q.Execute(workFunc); err != nil {
		return err
	}

	condition.Value = value
	condition.IsSubquery = false
	return nil
}
