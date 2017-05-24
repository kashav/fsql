package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kshvmdn/fsql/parser"
)

const (
	version = "0.1.0"

	uBYTE     = 1.0
	uKILOBYTE = 1024 * uBYTE
	uMEGABYTE = 1024 * uKILOBYTE
	uGIGABYTE = 1024 * uMEGABYTE
)

// Used to track which paths we've seen to avoid revisiting a directory.
var seen = make(map[string]bool, 0)

// Read the command line arguments for the query.
func readFlags() string {
	flag.Usage = func() {
		fmt.Printf("usage: %s [options] query\n", os.Args[0])
		flag.PrintDefaults()
	}

	versionPtr := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionPtr {
		fmt.Printf("fsql v%v\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(flag.Args()) > 1 {
		return strings.Join(flag.Args(), " ")
	}

	return flag.Args()[0]
}

func main() {
	input := readFlags()

	q, err := parser.Run(input)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.Fatal("Unexpected end of line")
		}
		log.Fatal(err)
	}

	q.Execute(func(path string, info os.FileInfo) {
		results := q.ApplyModifiers(path, info)

		if q.HasAttribute("mode") {
			fmt.Printf("%s", results["mode"])
			if q.HasAttribute("size", "time", "name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("size") {
			fmt.Printf("%v", results["size"])
			if q.HasAttribute("time", "name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("time") {
			fmt.Printf("%s", results["time"])
			if q.HasAttribute("name") {
				fmt.Print("\t")
			}
		}

		if q.HasAttribute("name") {
			fmt.Printf("%s", results["name"])
		}
		fmt.Printf("\n")
	})
}
