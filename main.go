package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kshvmdn/fsql/parser"
)

// Read the command line arguments for the query.
func readInput() string {
	if len(os.Args) == 1 {
		log.Fatal("Expected query.")
	}

	if len(os.Args) > 2 {
		return strings.Join(os.Args[1:], " ")
	}

	return os.Args[1]
}

func main() {
	input := readInput()

	q, err := parser.Run(input)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.Fatal("Unexpected end of line")
		}
		log.Fatal(err)
	}

	q.Execute(func(path string, info os.FileInfo) {
		if q.HasAttribute("mode") {
			fmt.Printf("%s", info.Mode())
			if q.HasAttribute("size", "time", "name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("size") {
			fmt.Printf("%d", info.Size())
			if q.HasAttribute("time", "name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("time") {
			fmt.Printf("%s", info.ModTime().Format(time.Stamp))
			if q.HasAttribute("name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("name") {
			fmt.Printf("%s", path)
		}

		fmt.Printf("\n")
	})
}
