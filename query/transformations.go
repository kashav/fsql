package query

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Function mapping of transformations
var allTransformations = map[string]func(interface{}, []string) interface{}{
	"format":   format,
	"upper":    upper,
	"fullpath": fullpath,
}

//format function
func format(inp interface{}, args []string) interface{} {

	if size, ok := inp.(int64); ok {
		if args[0] == "kb" {
			return fmt.Sprintf("%fkb", float64(size)/(1<<10))
		} else if args[0] == "mb" {
			return fmt.Sprintf("%fmb", float64(size)/(1<<20))
		} else if args[0] == "gb" {
			return fmt.Sprintf("%fgb", float64(size)/(1<<30))
		}
	} else {
		fmt.Fprintf(os.Stderr, "type for attribute:%s for function: format is not implemented\n", inp)
		os.Exit(1)
	}
	return ""
}

//Upper Function
func upper(inp interface{}, args []string) interface{} {
	if str, ok := inp.(string); ok {
		return strings.ToUpper(str)
	}
	fmt.Fprintf(os.Stderr, "type for attribute:%s for function: format is not implemented\n", inp)
	os.Exit(1)
	return ""
}

//Fullpath function
func fullpath(inp interface{}, args []string) interface{} {
	if str, ok := inp.(string); ok {
		res, _ := filepath.Abs(str)
		return res
	}
	fmt.Fprintf(os.Stderr, "type for attribute:%s for function: format is not implemented\n", inp)
	os.Exit(1)
	return ""
}

//PerformTransformations fetchs all transformation to be applied on a given attribute and applies them one by one
func (q *Query) PerformTransformations(attribute string, in interface{}) string {
	functions, ok := q.Transformations[attribute]
	res := in
	if ok {
		for _, function := range functions {
			if funcImpl, ok := allTransformations[function.Name]; ok {
				res = funcImpl(res, function.Arguments)
			} else {
				fmt.Fprintf(os.Stderr, "Given function:%s is not implemented\n", function.Name)
				os.Exit(1)
			}
		}
	}

	if v, ok := res.(int64); ok {
		return fmt.Sprintf("%d", v)
	} else if v, ok := res.(string); ok {
		return v
	}
	return ""
}
