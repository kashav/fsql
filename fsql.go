package fsql

import (
	"fmt"
	"os"

	"github.com/kshvmdn/fsql/parser"
	"github.com/kshvmdn/fsql/prompt"
	"github.com/kshvmdn/fsql/query"
)

// Version holds the current version number.
const Version = "0.2.1"

var q *query.Query

// Should be noted that this slice, is temporaly related to the order of printing
var attrs = [5]string{"mode", "size", "time", "hash", "name"}

// output prints the result value for each SELECTed attribute. Order is based
// on the order the attributes appear in attrs.
func output(path string, info os.FileInfo, result map[string]interface{}) {
	for i, attr := range attrs {
		if q.HasAttribute(attr) {
			fmt.Printf("%v", result[attr])
			if q.HasAttribute(attrs[i+1:]...) {
				fmt.Print("\t")
			}
		}
	}
	fmt.Print("\n")
}

// Run parses the input and executes the resultant query.
func Run(input string) (err error) {
	if q, err = parser.Run(input); err != nil {
		return err
	}
	if err = q.Execute(output); err != nil {
		return err
	}
	return nil
}

// RunInteractive starts the prompt and continuously calls Run until the
// process is exited or prompt.Run reads nothing.
func RunInteractive() error {
	for {
		if input := prompt.Run(); input == nil {
			return nil
		} else if err := Run(*input); err != nil {
			return err
		}
	}
}
