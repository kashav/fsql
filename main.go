package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cmp "./compare"
	"./query"
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

// Runs the appropriate cmp method for the provided condition.
func compare(condition query.Condition, file os.FileInfo) bool {
	var retval bool

	switch condition.Attribute {
	case "name":
		retval = cmp.Alpha(condition.Comparator, file.Name(), condition.Value)

	case "size":
		mult := uBYTE

		if len(condition.Value) > 2 {
			unit := strings.ToLower(condition.Value[len(condition.Value)-2:])
			switch unit {
			case "kb":
				mult = uKILOBYTE
			case "mb":
				mult = uMEGABYTE
			case "gb":
				mult = uGIGABYTE
			}

			if mult > 1 {
				condition.Value = condition.Value[:len(condition.Value)-2]
			}
		}

		size, err := strconv.ParseFloat(condition.Value, 64)
		if err != nil {
			return false
		}
		retval = cmp.Numeric(condition.Comparator, file.Size(), int64(size*mult))

	case "time":
		t, err := time.Parse("Jan 02 2006 15 04", condition.Value)
		if err != nil {
			return false
		}
		retval = cmp.Time(condition.Comparator, file.ModTime(), t)

	case "file":
		retval = cmp.File(condition.Comparator, file, condition.Value)
	}

	if condition.Negate {
		return !retval
	}

	return retval
}

// Return true iff path contains a substring of any element of exclusions.
func containsAny(exclusions []string, path string) bool {
	for _, exclusion := range exclusions {
		if strings.Contains(path, exclusion) {
			return true
		}
	}

	return false
}

func main() {
	input := readFlags()

	q, err := query.RunParser(input)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.Fatal("Unexpected end of line")
		}
		log.Fatal(err)
	}

	for _, src := range q.Sources["include"] {
		filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if path == "." || path == ".." || err != nil {
				return nil
			}

			if _, ok := seen[path]; ok {
				return nil
			}
			seen[path] = true

			// If this path is excluded or the condition is false, return.
			if containsAny(q.Sources["exclude"], path) ||
				!q.ConditionTree.Evaluate(info, compare) {
				return nil
			}

			if q.HasAttribute("mode") {
				fmt.Printf("%s", q.PerformTransformations("mode",info.Mode()))
				if q.HasAttribute("size", "time", "name") {
					fmt.Print("\t")
				}
			}

			if q.HasAttribute("size") {
				fmt.Printf("%s", q.PerformTransformations("size",info.Size()))
				if q.HasAttribute("time", "name") {
					fmt.Print("\t")
				}
			}

			if q.HasAttribute("time") {
				fmt.Printf("%s", q.PerformTransformations("time",info.ModTime().Format(time.Stamp)) )
				if q.HasAttribute("name") {
					fmt.Print("\t")
				}
			}

			if q.HasAttribute("name") {
				// TODO: Only show file name, instead of the full path?
				fmt.Printf("%s", q.PerformTransformations("name",path))
			}

			fmt.Printf("\n")
			return nil
		})
	}
}
