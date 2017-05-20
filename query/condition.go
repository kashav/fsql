package query

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
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
		return fmt.Sprintf("(%v)", root.Condition.String())
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
	Attribute  string
	Comparator tokenizer.TokenType
	Value      interface{}
	Negate     bool
	Subquery   *Query
	IsSubquery bool
}

func (c *Condition) String() string {
	return fmt.Sprintf(
		"{attribute: %s, comparator: %s, value: %v, negate: %t, subquery: %t}",
		c.Attribute, c.Comparator, c.Value, c.Negate, c.IsSubquery)
}

// Evaluate this single condition.
func (c *Condition) evaluate(path string, file os.FileInfo) bool {
	var retval bool

	switch c.Attribute {
	case "name":
		if !c.IsSubquery {
			retval = cmpAlpha(c.Comparator, file.Name(), c.Value)
		}
	case "path":
		if !c.IsSubquery {
			retval = cmpAlpha(c.Comparator, path, c.Value)
		}
	case "size":
		if _, ok := c.Value.(map[interface{}]bool); !ok {
			sizeStr := c.Value.(string)

			var multiplier float64
			if len(sizeStr) > 2 {
				switch strings.ToLower(sizeStr[len(sizeStr)-2:]) {
				case "kb":
					multiplier = uKILOBYTE
				case "mb":
					multiplier = uMEGABYTE
				case "gb":
					multiplier = uGIGABYTE
				default:
					multiplier = uBYTE
				}

				if multiplier != uBYTE {
					sizeStr = sizeStr[:len(sizeStr)-2]
				}
			}

			size, err := strconv.ParseFloat(sizeStr, 64)
			if err != nil {
				return false
			}
			retval = cmpNumeric(c.Comparator, file.Size(), int64(size*multiplier))
		} else {
			retval = cmpNumeric(c.Comparator, file.Size(), c.Value)
		}
	case "time":
		if _, ok := c.Value.(map[interface{}]bool); !ok {
			t, err := time.Parse("Jan 02 2006 15 04", c.Value.(string))
			if err != nil {
				return false
			}
			retval = cmpTime(c.Comparator, file.ModTime(), t)
		} else {
			retval = cmpTime(c.Comparator, file.ModTime(), c.Value)
		}
	case "file":
		retval = cmpFile(c.Comparator, file, c.Value)
	}

	if c.Negate {
		return !retval
	}

	return retval
}
