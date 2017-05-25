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

// Run pre-order traversal on the ConditionNode tree rooted at root and
// evaluates each conditional along the path with the provided compare method.
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

		if _, ok := root.Condition.Value.(map[interface{}]bool); ok {
			// Array of values, returned from evaluating a subquery.
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

// Evaluate this Condition.
func (c *Condition) evaluate(path string, file os.FileInfo) bool {
	var retval bool

	switch c.Attribute {
	case "name":
		retval = evalAlpha(c.Comparator, file.Name(), c.Value.(string))

	case "size":
		var size int64

		switch c.Value.(type) {
		case float64:
			size = int64(c.Value.(float64))
		case string:
			sizeFloat, err := strconv.ParseFloat(c.Value.(string), 10)
			if err != nil {
				log.Fatalln(err)
			}
			size = int64(sizeFloat)
		}

		retval = evalNumeric(c.Comparator, file.Size(), size)

	case "time":
		switch c.Value.(type) {
		case time.Time:
		case string:
			t, err := time.Parse("Jan 02 2006 15 04", c.Value.(string))
			if err != nil {
				log.Fatalln(err)
			}
			c.Value = t
		}

		retval = evalTime(c.Comparator, file.ModTime(), c.Value.(time.Time))

	case "file":
		retval = evalFile(c.Comparator, file, c.Value)
	}

	if c.Negate {
		return !retval
	}

	return retval
}
