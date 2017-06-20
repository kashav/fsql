package query

import (
	"errors"
	"fmt"
	"os"

	"github.com/kshvmdn/fsql/evaluate"
	"github.com/kshvmdn/fsql/tokenizer"
	"github.com/kshvmdn/fsql/transform"
)

// ConditionNode represents a single node of a query's WHERE clause tree.
type ConditionNode struct {
	Type      *tokenizer.TokenType
	Left      *ConditionNode
	Right     *ConditionNode
	Condition *Condition
}

func (root *ConditionNode) String() string {
	if root == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{%v (%v %v) %v}", root.Type, root.Left, root.Right,
		root.Condition)
}

// evaluateTree runs pre-order traversal on the ConditionNode tree rooted at
// root and evaluates each conditional along the path with the provided compare
// method.
func (root *ConditionNode) evaluateTree(path string, info os.FileInfo) (bool, error) {
	if root == nil {
		return true, nil
	}

	if root.Condition != nil {
		if root.Condition.IsSubquery {
			// Unevaluated subquery.
			// TODO: Handle this case.
			return false, errors.New("not implemented")
		}

		if !root.Condition.Parsed {
			if err := root.Condition.applyModifiers(); err != nil {
				return false, err
			}
		}

		return root.Condition.evaluate(path, info)
	}

	if *root.Type == tokenizer.And {
		if ok, err := root.Left.evaluateTree(path, info); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
		return root.Right.evaluateTree(path, info)
	}

	if *root.Type == tokenizer.Or {
		if ok, err := root.Left.evaluateTree(path, info); err != nil {
			return false, nil
		} else if ok {
			return true, nil
		}
		return root.Right.evaluateTree(path, info)
	}

	return false, nil
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute          string
	AttributeModifiers []Modifier
	Parsed             bool

	Operator tokenizer.TokenType
	Value    interface{}
	Negate   bool

	Subquery   *Query
	IsSubquery bool
}

// ApplyModifiers applies each modifier to the value of this Condition.
func (c *Condition) applyModifiers() error {
	value := c.Value

	for _, m := range c.AttributeModifiers {
		var err error
		value, err = transform.Parse(&transform.ParseParams{
			Attribute: c.Attribute,
			Value:     value,
			Name:      m.Name,
			Args:      m.Arguments,
		})
		if err != nil {
			return err
		}
	}

	c.Value = value
	c.Parsed = true
	return nil
}

// evaluate runs the respective evaluate function for this Condition.
func (c *Condition) evaluate(path string, file os.FileInfo) (bool, error) {
	// FIXME: This is a bit of a hack. We can't pass c.AttributeModifiers, since
	// that'll cause a import cycle, so we have to recreate the attribute
	// modifiers slice using a separate type defined in evaluate.
	modifiers := make([]evaluate.Modifier, len(c.AttributeModifiers))
	for i, m := range c.AttributeModifiers {
		modifiers[i] = evaluate.Modifier{Name: m.Name, Arguments: m.Arguments}
	}

	o := &evaluate.Opts{
		Path:      path,
		File:      file,
		Attribute: c.Attribute,
		Modifiers: modifiers,
		Operator:  c.Operator,
		Value:     c.Value,
	}
	result, err := evaluate.Evaluate(o)
	if err != nil {
		return false, err
	}
	if c.Negate {
		return !result, nil
	}
	return result, nil
}
