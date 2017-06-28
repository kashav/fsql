package fsql

import (
	"fmt"
	"os"

	"github.com/kshvmdn/fsql/parser"
)

// Run parses the input and executes the resultant query.
func Run(input string) error {
	q, err := parser.Run(input)
	if err != nil {
		return err
	}

	// Find length of the longest name to normalize name output.
	var max = 0
	var results = make([]map[string]interface{}, 0)

	err = q.Execute(
		func(path string, info os.FileInfo, result map[string]interface{}) {
			results = append(results, result)
			if !q.HasAttribute("name") {
				return
			}
			if s, ok := result["name"].(string); ok && len(s) > max {
				max = len(s)
			}
		},
	)
	if err != nil {
		return err
	}

	for _, result := range results {
		for j, attribute := range q.Attributes {
			// If the current attribute is "name", pad the output string by `max`
			// spaces.
			format := "%v"
			if attribute == "name" {
				format = fmt.Sprintf("%%-%ds", max)
			}
			fmt.Printf(format, result[attribute])
			if j != len(q.Attributes)-1 {
				fmt.Print("\t")
			}
		}
		fmt.Print("\n")
	}

	return nil
}
