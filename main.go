package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	cmp "github.com/kshvmdn/fsql/compare"
	"github.com/kshvmdn/fsql/query"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
)

func compare(condition query.Condition, file os.FileInfo) bool {
	var retval bool

	switch condition.Attribute {
	case "name":
		retval = cmp.Alpha(condition.Comparator, file.Name(), condition.Value)

	case "size":
		unit := strings.ToLower(condition.Value[len(condition.Value)-2:])
		mult := BYTE
		switch unit {
		case "kb":
			mult = KILOBYTE
		case "mb":
			mult = MEGABYTE
		case "gb":
			mult = GIGABYTE
		}

		if mult > 1 {
			condition.Value = condition.Value[:len(condition.Value)-2]
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

	return (condition.Negate && !retval) || retval
}

func containsAny(path string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if strings.Contains(path, exclusion) {
			return true
		}
	}

	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected query.")
		os.Exit(1)
	}
	input := os.Args[1]

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	q, err := query.Run(input)

	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.Fatal("Unexpected end of line.")
		}

		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, src := range q.Sources["include"] {
		wg.Add(1)

		go func(src string) {
			defer wg.Done()

			if strings.Contains(src, "~") {
				src = filepath.Join(usr.HomeDir, src[1:])
			}

			filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
				if path == "." || path == ".." || containsAny(path, q.Sources["exclude"]) {
					return nil
				}

				if q.ConditionTree.Evaluate(info, compare) {
					if q.HasAttribute("mode") {
						fmt.Printf("%s\t", info.Mode())
					}

					if q.HasAttribute("size") {
						fmt.Printf("%d\t", info.Size())
					}

					if q.HasAttribute("time") {
						fmt.Printf("%s\t", info.ModTime().Format(time.Stamp))
					}

					if q.HasAttribute("name") {
						if strings.Contains(path, usr.HomeDir) {
							path = filepath.Join("~", path[len(usr.HomeDir):])
						}
						fmt.Printf("%s", path)
					}

					fmt.Printf("\n")
				}

				return nil
			})
		}(src)
	}

	wg.Wait()
}
