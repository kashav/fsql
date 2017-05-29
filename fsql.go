package fsql

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kshvmdn/fsql/parser"
	"github.com/kshvmdn/fsql/prompt"
	"github.com/kshvmdn/fsql/query"
)

// Version holds the current version number.
const Version = "0.1.0"

var q *query.Query

// output prints the result value for each SELECTed output.
func output(path string, info os.FileInfo, result map[string]interface{}) {
	if q.HasAttribute("mode") {
		fmt.Printf("%s", result["mode"])
		if q.HasAttribute("size", "time", "name") {
			fmt.Print("\t")
		}
	}
	if q.HasAttribute("size") {
		fmt.Printf("%v", result["size"])
		if q.HasAttribute("time", "name") {
			fmt.Print("\t")
		}
	}
	if q.HasAttribute("time") {
		fmt.Printf("%s", result["time"])
		if q.HasAttribute("name") {
			fmt.Print("\t")
		}
	}
	if q.HasAttribute("name") {
		fmt.Printf("%s", result["name"])
	}
	fmt.Printf("\n")
}

// Run parses the input and executes the resultant query.
func Run(input string) (err error) {
	q, err = parser.Run(input)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return errors.New("unexpected end of line")
		}
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
