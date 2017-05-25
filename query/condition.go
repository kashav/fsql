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
	Type      tokenizer.TokenType
	Condition *Condition
	Left      *ConditionNode
	Right     *ConditionNode
}

func (root *ConditionNode) String() string {
	if root == nil {
		return "nil"
	}

	if root.Condition != nil {
		return fmt.Sprintf("(%v)", root.Condition)
	}

	return fmt.Sprintf("(%s (%s, %s))", root.Type, root.Left.String(),
		root.Right.String())
}

// evaluateTree runs pre-order traversal on the ConditionNode tree rooted at
// root and evaluates each conditional along the path with the provided compare
// method.
func (root *ConditionNode) evaluateTree(path string, file os.FileInfo) bool {
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

		return root.Condition.evaluate(path, file)
	}

	if root.Type == tokenizer.And {
		return root.Left.evaluateTree(path, file) &&
			root.Right.evaluateTree(path, file)
	}

	if root.Type == tokenizer.Or {
		if root.Left.evaluateTree(path, file) {
			return true
		}
		return root.Right.evaluateTree(path, file)
	}

	return false
}

// Condition represents a WHERE condition.
type Condition struct {
	Attribute          string
	AttributeModifiers []Modifier
	Parsed             bool

	Comparator tokenizer.TokenType
	Value      interface{}
	Negate     bool

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
			Value:     c.Value,
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
		return cmpAlpha(c.Comparator, file.Name(), c.Value.(string))

	case []string:
		return cmpAlpha(c.Comparator, file.Name(), c.Value.([]string))

	case map[interface{}]bool:
		return cmpAlpha(c.Comparator, file.Name(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateSize evaluates a Condition with attribute `size`.
func (c *Condition) evaluateSize(path string, file os.FileInfo) bool {
	switch c.Value.(type) {
	case float64:
		return cmpNumeric(c.Comparator, file.Size(), int64(c.Value.(float64)))

	case string:
		size, err := strconv.ParseFloat(c.Value.(string), 10)
		if err != nil {
			log.Fatalln(err)
		}

		return cmpNumeric(c.Comparator, file.Size(), int64(size))

	case map[interface{}]bool:
		return cmpNumeric(c.Comparator, file.Size(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateTime evaluates a Condition with attribute `time`.
func (c *Condition) evaluateTime(path string, file os.FileInfo) bool {
	switch c.Value.(type) {
	case string:
		t, err := time.Parse("Jan 02 2006 15 04", c.Value.(string))
		if err != nil {
			log.Fatalln(err)
		}
		return cmpTime(c.Comparator, file.ModTime(), t)

	case time.Time:
		return cmpTime(c.Comparator, file.ModTime(), c.Value.(time.Time))

	case map[interface{}]bool:
		return cmpTime(c.Comparator, file.ModTime(), c.Value.(map[interface{}]bool))
	}

	return false
}

// evaluateMode evaluates a Condition with attribute `mode`.
func (c *Condition) evaluateMode(path string, file os.FileInfo) bool {
	return cmpMode(c.Comparator, file, c.Value)
}
