package query

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
func (root *ConditionNode) evaluateTree(path string, info os.FileInfo) bool {
	if root == nil {
		return true
	}

	if root.Condition != nil {
		if root.Condition.IsSubquery {
			// Unevaluated subquery.
			// TODO: Handle this case.
			return false
		}

		if !root.Condition.Parsed {
			if err := root.Condition.applyModifiers(); err != nil {
				log.Fatal(err.Error())
			}
		}

		return root.Condition.evaluate(path, info)
	}

	if *root.Type == tokenizer.And {
		if !root.Left.evaluateTree(path, info) {
			return false
		}
		return root.Right.evaluateTree(path, info)
	}

	if *root.Type == tokenizer.Or {
		if root.Left.evaluateTree(path, info) {
			return true
		}
		return root.Right.evaluateTree(path, info)
	}

	return false
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
func (c *Condition) evaluate(path string, file os.FileInfo) bool {
	var retval bool

	switch c.Attribute {
	case "name":
		retval = c.evaluateName(path, file)
	case "size":
		retval = c.evaluateSize(path, file)
	case "time":
		retval = c.evaluateTime(path, file)
	case "mode":
		retval = c.evaluateMode(path, file)
	}

	if c.Negate {
		return !retval
	}

	return retval
}

// evaluateName evaluates a Condition with attribute `name`.
func (c *Condition) evaluateName(path string, file os.FileInfo) bool {
	switch c.Value.(type) {
	case string:
		return cmpAlpha(c.Operator, file.Name(), c.Value.(string))
	case []string:
		return cmpAlpha(c.Operator, file.Name(), c.Value.([]string))
	case map[interface{}]bool:
		return cmpAlpha(c.Operator, file.Name(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateSize evaluates a Condition with attribute `size`.
func (c *Condition) evaluateSize(path string, file os.FileInfo) bool {
	switch c.Value.(type) {
	case float64:
		return cmpNumeric(c.Operator, file.Size(), int64(c.Value.(float64)))
	case string:
		size, err := strconv.ParseFloat(c.Value.(string), 10)
		if err != nil {
			log.Fatal(err.Error())
		}
		return cmpNumeric(c.Operator, file.Size(), int64(size))
	case map[interface{}]bool:
		return cmpNumeric(c.Operator, file.Size(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateTime evaluates a Condition with attribute `time`.
func (c *Condition) evaluateTime(path string, file os.FileInfo) bool {
	switch c.Value.(type) {
	case string:
		t, err := time.Parse("Jan 02 2006 15 04", c.Value.(string))
		if err != nil {
			log.Fatal(err.Error())
		}
		return cmpTime(c.Operator, file.ModTime(), t)
	case time.Time:
		return cmpTime(c.Operator, file.ModTime(), c.Value.(time.Time))
	case map[interface{}]bool:
		return cmpTime(c.Operator, file.ModTime(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateMode evaluates a Condition with attribute `mode`.
func (c *Condition) evaluateMode(path string, file os.FileInfo) bool {
	return cmpMode(c.Operator, file, c.Value)
}
